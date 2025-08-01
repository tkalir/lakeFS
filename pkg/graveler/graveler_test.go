package graveler_test

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/pkg/catalog"
	"github.com/treeverse/lakefs/pkg/catalog/testutils"
	"github.com/treeverse/lakefs/pkg/graveler"
	"github.com/treeverse/lakefs/pkg/graveler/mock"
	"github.com/treeverse/lakefs/pkg/graveler/testutil"
	"github.com/treeverse/lakefs/pkg/kv"
)

type Hooks struct {
	Called           []string
	Errs             map[string]error
	RunID            string
	RepositoryID     graveler.RepositoryID
	StorageNamespace graveler.StorageNamespace
	BranchID         graveler.BranchID
	SourceRef        graveler.Ref
	CommitID         graveler.CommitID
	Commit           graveler.Commit
	TagID            graveler.TagID
}

var ErrGravelerUpdate = errors.New("test update error")

func (h *Hooks) PrepareCommitHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PrepareCommitHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.StorageNamespace = record.Repository.StorageNamespace
	h.BranchID = record.BranchID
	return h.Errs["PrepareCommitHook"]
}

func (h *Hooks) PreCommitHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreCommitHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.StorageNamespace = record.Repository.StorageNamespace
	h.BranchID = record.BranchID
	h.Commit = record.Commit
	return h.Errs["PreCommitHook"]
}

func (h *Hooks) PostCommitHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PostCommitHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
	h.CommitID = record.CommitID
	h.Commit = record.Commit
	return h.Errs["PostCommitHook"]
}

func (h *Hooks) PreMergeHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreMergeHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.StorageNamespace = record.Repository.StorageNamespace
	h.BranchID = record.BranchID
	h.SourceRef = record.SourceRef
	h.Commit = record.Commit
	return h.Errs["PreMergeHook"]
}

func (h *Hooks) PostMergeHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PostMergeHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.StorageNamespace = record.Repository.StorageNamespace
	h.BranchID = record.BranchID
	h.SourceRef = record.SourceRef
	h.CommitID = record.CommitID
	h.Commit = record.Commit
	return h.Errs["PostMergeHook"]
}

func (h *Hooks) PreCreateTagHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreCreateTagHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.CommitID = record.CommitID
	h.TagID = record.TagID
	return h.Errs["PreCreateTagHook"]
}

func (h *Hooks) PostCreateTagHook(_ context.Context, record graveler.HookRecord) {
	h.Called = append(h.Called, "PostCreateTagHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.CommitID = record.CommitID
	h.TagID = record.TagID
}

func (h *Hooks) PreDeleteTagHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreDeleteTagHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.TagID = record.TagID
	return h.Errs["PreDeleteTagHook"]
}

func (h *Hooks) PostDeleteTagHook(_ context.Context, record graveler.HookRecord) {
	h.Called = append(h.Called, "PostDeleteTagHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.TagID = record.TagID
}

func (h *Hooks) PreCreateBranchHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreCreateBranchHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
	h.CommitID = record.CommitID
	h.SourceRef = record.SourceRef
	return h.Errs["PreCreateBranchHook"]
}

func (h *Hooks) PostCreateBranchHook(_ context.Context, record graveler.HookRecord) {
	h.Called = append(h.Called, "PostCreateBranchHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
	h.CommitID = record.CommitID
	h.SourceRef = record.SourceRef
}

func (h *Hooks) PreDeleteBranchHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreDeleteBranchHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
	return h.Errs["PreDeleteBranchHook"]
}

func (h *Hooks) PostDeleteBranchHook(_ context.Context, record graveler.HookRecord) {
	h.Called = append(h.Called, "PostDeleteBranchHook")
	h.StorageNamespace = record.Repository.StorageNamespace
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
}

func (h *Hooks) PreRevertHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreRevertHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.StorageNamespace = record.Repository.StorageNamespace
	h.BranchID = record.BranchID
	h.Commit = record.Commit
	return h.Errs["PreRevertHook"]
}

func (h *Hooks) PostRevertHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PostRevertHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
	h.CommitID = record.CommitID
	h.Commit = record.Commit
	return h.Errs["PostRevertHook"]
}

func (h *Hooks) PreCherryPickHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PreCherryPickHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.StorageNamespace = record.Repository.StorageNamespace
	h.BranchID = record.BranchID
	h.Commit = record.Commit
	return h.Errs["PreCherryPickHook"]
}

func (h *Hooks) PostCherryPickHook(_ context.Context, record graveler.HookRecord) error {
	h.Called = append(h.Called, "PostCherryPickHook")
	h.RepositoryID = record.Repository.RepositoryID
	h.BranchID = record.BranchID
	h.CommitID = record.CommitID
	h.Commit = record.Commit
	return h.Errs["PostCherryPickHook"]
}

func (h *Hooks) NewRunID() string {
	return ""
}

func newGraveler(t *testing.T, committedManager graveler.CommittedManager, stagingManager graveler.StagingManager,
	refManager graveler.RefManager, gcManager graveler.GarbageCollectionManager,
	protectedBranchesManager graveler.ProtectedBranchesManager,
) catalog.Store {
	t.Helper()

	return graveler.NewGraveler(committedManager, stagingManager, refManager, gcManager, protectedBranchesManager, nil)
}

func TestGraveler_List(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		r           catalog.Store
		expectedErr error
		expected    []*graveler.ValueRecord
	}{
		{
			name: "one committed one staged no paths",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("bar"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("bar"), Value: &graveler.Value{}}, {Key: graveler.Key("foo"), Value: &graveler.Value{}}},
		},
		{
			name: "one compacted one staged no paths",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{}}}), MetaRangeID: "mr1"},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("bar"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mr1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("bar"), Value: &graveler.Value{}}, {Key: graveler.Key("foo"), Value: &graveler.Value{}}},
		},
		{
			name: "same path different file",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("original")}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("other")}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("other")}}},
		},
		{
			name: "same path different file compacted",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("original")}}}), MetaRangeID: "mr1"},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("other")}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mr1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("foo"), Value: &graveler.Value{Identity: []byte("other")}}},
		},
		{
			name: "one committed one staged no paths - with prefix",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("prefix/foo"), Value: &graveler.Value{}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("prefix/bar"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("prefix/bar"), Value: &graveler.Value{}}, {Key: graveler.Key("prefix/foo"), Value: &graveler.Value{}}},
		},
		{
			name: "one compacted one staged no paths - with prefix",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("prefix/foo"), Value: &graveler.Value{}}}), MetaRangeID: "mr1"},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("prefix/bar"), Value: &graveler.Value{}}})},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mr1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expected: []*graveler.ValueRecord{{Key: graveler.Key("prefix/bar"), Value: &graveler.Value{}}, {Key: graveler.Key("prefix/foo"), Value: &graveler.Value{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listing, err := tt.r.List(ctx, repository, "", 0)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("wrong error, expected:%v got:%v", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expectedErr
			}
			defer listing.Close()
			var recs []*graveler.ValueRecord
			for listing.Next() {
				recs = append(recs, listing.Value())
			}
			if diff := deep.Equal(recs, tt.expected); diff != nil {
				t.Fatal("List() diff found", diff)
			}
		})
	}
}

func TestGraveler_Get(t *testing.T) {
	errTest := errors.New("some kind of err")
	tests := []struct {
		name                string
		r                   catalog.Store
		expectedValueResult graveler.Value
		expectedErr         error
	}{
		{
			name: "commit - exists",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, nil,
				&testutil.RefsFake{RefType: graveler.ReferenceTypeCommit, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("committed")},
		},
		{
			name: "commit - not found",
			r: newGraveler(t, &testutil.CommittedFake{Err: graveler.ErrNotFound}, nil,
				&testutil.RefsFake{RefType: graveler.ReferenceTypeCommit, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			), expectedErr: graveler.ErrNotFound,
		},
		{
			name: "commit - error",
			r: newGraveler(t, &testutil.CommittedFake{Err: errTest}, nil,
				&testutil.RefsFake{RefType: graveler.ReferenceTypeCommit, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			), expectedErr: errTest,
		},
		{
			name: "branch - only staged",
			r: newGraveler(t, &testutil.CommittedFake{Err: graveler.ErrNotFound}, &testutil.StagingFake{Value: &graveler.Value{Identity: []byte("staged")}},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("staged")},
		},
		{
			name: "branch - committed and staged",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, &testutil.StagingFake{Value: &graveler.Value{Identity: []byte("staged")}},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("staged")},
		},
		{
			name: "branch - only committed",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, &testutil.StagingFake{Err: graveler.ErrNotFound},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, StagingToken: "token"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("committed")},
		},
		{
			name: "branch - tombstone",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("committed")}}}, &testutil.StagingFake{Value: nil},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedErr: graveler.ErrNotFound,
		},
		{
			name: "branch - staged return error",
			r: newGraveler(t, &testutil.CommittedFake{}, &testutil.StagingFake{Err: errTest},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedErr: errTest,
		},
		{
			name: "branch - only compacted",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("compacted")}}, MetaRangeID: "mir1"}, &testutil.StagingFake{Err: graveler.ErrNotFound},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mir1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("compacted")},
		},
		{
			name: "branch - staged and compacted",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("compacted")}}, MetaRangeID: "mir1"}, &testutil.StagingFake{Value: &graveler.Value{Identity: []byte("staged")}},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mir1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("staged")},
		},
		{
			name: "branch - deleted from staged, exists in compaction",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"key": {Identity: []byte("compacted")}}, MetaRangeID: "mir1"}, &testutil.StagingFake{Value: nil},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mir1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedErr: graveler.ErrNotFound,
		},
		{
			name: "branch - exists in staging and not in compaction",
			r: newGraveler(t, &testutil.CommittedFake{MetaRangeID: "mir1"}, &testutil.StagingFake{Value: &graveler.Value{Identity: []byte("staged")}},
				&testutil.RefsFake{RefType: graveler.ReferenceTypeBranch, StagingToken: "token1", Commits: map[graveler.CommitID]*graveler.Commit{"": {}}, BaseMetaRangeID: "mir1"}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedValueResult: graveler.Value{Identity: []byte("staged")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Value, err := tt.r.Get(context.Background(), repository, "", []byte("key"))
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("wrong error, expected:%v got:%v", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expected error
			}
			if string(tt.expectedValueResult.Identity) != string(Value.Identity) {
				t.Errorf("wrong Value address, expected:%s got:%s", tt.expectedValueResult.Identity, Value.Identity)
			}
		})
	}
}

func TestGraveler_Set(t *testing.T) {
	newSetVal := &graveler.ValueRecord{Key: []byte("key"), Value: &graveler.Value{Data: []byte("newValue"), Identity: []byte("newIdentity")}}
	sampleVal := &graveler.Value{Identity: []byte("sampleIdentity"), Data: []byte("sampleValue")}
	tests := []struct {
		name                string
		ifAbsent            bool
		expectedValueResult *graveler.ValueRecord
		expectedErr         error
		committedMgr        *testutil.CommittedFake
		stagingMgr          *testutil.StagingFake
		refMgr              *testutil.RefsFake
	}{
		{
			name:                "with nothing before",
			committedMgr:        &testutil.CommittedFake{},
			stagingMgr:          &testutil.StagingFake{},
			refMgr:              &testutil.RefsFake{Branch: &graveler.Branch{}},
			expectedValueResult: newSetVal,
		},
		{
			name:                "with committed key",
			committedMgr:        &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{string(newSetVal.Key): {Data: []byte("dsa"), Identity: []byte("asd")}}},
			stagingMgr:          &testutil.StagingFake{},
			refMgr:              &testutil.RefsFake{Branch: &graveler.Branch{CommitID: "commit1"}},
			expectedValueResult: newSetVal,
		},
		{
			name:                "overwrite no prior value",
			committedMgr:        &testutil.CommittedFake{Err: graveler.ErrNotFound},
			stagingMgr:          &testutil.StagingFake{},
			refMgr:              &testutil.RefsFake{Branch: &graveler.Branch{CommitID: "bla"}, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			expectedValueResult: newSetVal,
			ifAbsent:            true,
		},
		{
			name:                "overwrite with prior committed value",
			committedMgr:        &testutil.CommittedFake{},
			stagingMgr:          &testutil.StagingFake{},
			refMgr:              &testutil.RefsFake{Branch: &graveler.Branch{CommitID: "bla"}, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			expectedValueResult: nil,
			expectedErr:         graveler.ErrPreconditionFailed,
			ifAbsent:            true,
		},
		{
			name:                "overwrite with prior staging value",
			committedMgr:        &testutil.CommittedFake{},
			stagingMgr:          &testutil.StagingFake{Values: map[string]map[string]*graveler.Value{"st": {"key": sampleVal}}, LastSetValueRecord: &graveler.ValueRecord{Key: []byte("key"), Value: sampleVal}},
			refMgr:              &testutil.RefsFake{Branch: &graveler.Branch{CommitID: "bla", StagingToken: "st"}, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			expectedValueResult: &graveler.ValueRecord{Key: []byte("key"), Value: sampleVal},
			expectedErr:         graveler.ErrPreconditionFailed,
			ifAbsent:            true,
		},
		{
			name:                "overwrite with prior staging tombstone",
			committedMgr:        &testutil.CommittedFake{Err: graveler.ErrNotFound},
			stagingMgr:          &testutil.StagingFake{Values: map[string]map[string]*graveler.Value{"st1": {"key": nil}, "st2": {"key": sampleVal}}, LastSetValueRecord: &graveler.ValueRecord{Key: []byte("key"), Value: sampleVal}},
			refMgr:              &testutil.RefsFake{Branch: &graveler.Branch{CommitID: "bla", StagingToken: "st1", SealedTokens: []graveler.StagingToken{"st2"}}, Commits: map[graveler.CommitID]*graveler.Commit{"": {}}},
			expectedValueResult: newSetVal,
			ifAbsent:            true,
		},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newGraveler(t, tt.committedMgr, tt.stagingMgr, tt.refMgr, nil, testutil.NewProtectedBranchesManagerFake())
			err := store.Set(ctx, repository, "branch-1", newSetVal.Key, *newSetVal.Value, graveler.WithIfAbsent(tt.ifAbsent))
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Set() - error: %v, expected: %v", err, tt.expectedErr)
			}
			lastVal := tt.stagingMgr.LastSetValueRecord
			if err == nil {
				require.Equal(t, tt.expectedValueResult, lastVal)
			} else {
				require.NotEqual(t, &tt.expectedValueResult, lastVal)
			}
		})
	}
}

func TestGravelerSet_Advanced(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	committedMgr := &testutil.CommittedFake{}
	newSetVal := &graveler.ValueRecord{Key: []byte("key"), Value: &graveler.Value{Data: []byte("newValue"), Identity: []byte("newIdentity")}}
	// RefManager mock base setup
	refMgr := mock.NewMockRefManager(ctrl)
	refExpect := refMgr.EXPECT()
	refExpectCommitNotFound := func() {
		refExpect.ParseRef(gomock.Any()).Times(1).Return(graveler.RawRef{BaseRef: ""}, nil)
		refExpect.ResolveRawRef(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.ResolvedRef{
			BranchRecord: graveler.BranchRecord{
				Branch: &graveler.Branch{},
			},
		}, nil)
		refExpect.GetCommit(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, graveler.ErrCommitNotFound)
	}

	t.Run("update failure", func(t *testing.T) {
		stagingMgr := &testutil.StagingFake{
			UpdateErr: ErrGravelerUpdate,
		}
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{}, nil)
		refExpectCommitNotFound()
		store := newGraveler(t, committedMgr, stagingMgr, refMgr, nil, testutil.NewProtectedBranchesManagerFake())
		err := store.Set(ctx, repository, "branch-1", newSetVal.Key, *newSetVal.Value, graveler.WithIfAbsent(true))
		require.ErrorIs(t, err, ErrGravelerUpdate)
		require.Nil(t, stagingMgr.LastSetValueRecord)
	})

	t.Run("branch deleted after update", func(t *testing.T) {
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{}, nil)
		refExpectCommitNotFound()
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, graveler.ErrNotFound)
		stagingMgr := &testutil.StagingFake{}
		store := newGraveler(t, committedMgr, stagingMgr, refMgr, nil, testutil.NewProtectedBranchesManagerFake())
		err := store.Set(ctx, repository, "branch-1", newSetVal.Key, *newSetVal.Value, graveler.WithIfAbsent(true))
		require.ErrorIs(t, err, graveler.ErrNotFound)
		require.Equal(t, newSetVal, stagingMgr.LastSetValueRecord)
	})

	t.Run("branch token changed after update - one retry", func(t *testing.T) {
		// Test safeBranchWrite retries when token changed after update a single time and then succeeds
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{}, nil)
		refExpectCommitNotFound()
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(&graveler.Branch{
			StagingToken: "new_token",
		}, nil)
		refExpectCommitNotFound()
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
			StagingToken: "new_token",
		}, nil)
		stagingMgr := &testutil.StagingFake{}
		store := newGraveler(t, committedMgr, stagingMgr, refMgr, nil, testutil.NewProtectedBranchesManagerFake())
		err := store.Set(ctx, repository, "branch-1", newSetVal.Key, *newSetVal.Value, graveler.WithIfAbsent(true))
		require.Nil(t, err)
		require.Equal(t, newSetVal, stagingMgr.LastSetValueRecord)
	})

	t.Run("branch token changed max retries", func(t *testing.T) {
		// Test safeBranchWrite retries when token changed after update, reach the maximal number of retries and then succeed
		for i := 0; i < graveler.BranchWriteMaxTries-1; i++ {
			refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
				StagingToken: graveler.StagingToken("new_token_" + strconv.Itoa(i)),
			}, nil)
			refExpectCommitNotFound()
			refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
				StagingToken: graveler.StagingToken("new_token_" + strconv.Itoa(i+1)),
			}, nil)
		}
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
			StagingToken: graveler.StagingToken("new_token_" + strconv.Itoa(graveler.BranchWriteMaxTries)),
		}, nil)
		refExpectCommitNotFound()
		refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
			StagingToken: graveler.StagingToken("new_token_" + strconv.Itoa(graveler.BranchWriteMaxTries)),
		}, nil)
		stagingMgr := &testutil.StagingFake{}
		store := newGraveler(t, committedMgr, stagingMgr, refMgr, nil, testutil.NewProtectedBranchesManagerFake())
		err := store.Set(ctx, repository, "branch-1", newSetVal.Key, *newSetVal.Value, graveler.WithIfAbsent(true))
		require.Nil(t, err)
		require.Equal(t, newSetVal, stagingMgr.LastSetValueRecord)
	})

	t.Run("branch token changed retry exhausted", func(t *testing.T) {
		// Test safeBranchWrite retries when token changed after update, exceed the maximal number of retries and expect fail
		for i := 0; i < graveler.BranchWriteMaxTries; i++ {
			refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
				StagingToken: graveler.StagingToken("new_token_" + strconv.Itoa(i)),
			}, nil)
			refExpectCommitNotFound()
			refMgr.EXPECT().GetBranch(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&graveler.Branch{
				StagingToken: graveler.StagingToken("new_token_" + strconv.Itoa(i+1)),
			}, nil)
		}
		stagingMgr := &testutil.StagingFake{}
		store := newGraveler(t, committedMgr, stagingMgr, refMgr, nil, testutil.NewProtectedBranchesManagerFake())
		err := store.Set(ctx, repository, "branch-1", newSetVal.Key, *newSetVal.Value, graveler.WithIfAbsent(true))
		require.ErrorIs(t, err, graveler.ErrTooManyTries)
		require.Equal(t, newSetVal, stagingMgr.LastSetValueRecord)
	})
}

func TestGravelerGet_Advanced(t *testing.T) {
	tests := []struct {
		name                string
		r                   catalog.Store
		expectedValueResult graveler.Value
		expectedErr         error
	}{
		{
			name: "branch - staged with sealed tokens",
			r: newGraveler(t, &testutil.CommittedFake{ValuesByKey: map[string]*graveler.Value{"staged": {
				Identity: []byte("BAD"),
				Data:     nil,
			}}},
				&testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token1": {"staged": {
							Identity: []byte("stagedA"),
							Data:     nil,
						}},
						"token2": {"foo": {
							Identity: []byte("stagedB"),
							Data:     nil,
						}},
					},
				},
				&testutil.RefsFake{
					RefType:      graveler.ReferenceTypeBranch,
					StagingToken: "token1",
					SealedTokens: []graveler.StagingToken{"token2", "token3"},
					Commits:      map[graveler.CommitID]*graveler.Commit{"": {}},
				},
				nil, testutil.NewProtectedBranchesManagerFake()),
			expectedValueResult: graveler.Value{Identity: []byte("stagedA")},
		},
		{
			name: "branch - no staged with sealed tokens",
			r: newGraveler(t, &testutil.CommittedFake{Err: graveler.ErrNotFound},
				&testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token2": {"foo": {
							Identity: []byte("stagedZ"),
							Data:     nil,
						}},
						"token3": {"staged": {
							Identity: []byte("stagedA"),
							Data:     nil,
						}},
						"token4": {"staged": {
							Identity: []byte("stagedB"),
							Data:     nil,
						}},
					},
				},
				&testutil.RefsFake{
					RefType:      graveler.ReferenceTypeBranch,
					StagingToken: "token1",
					SealedTokens: []graveler.StagingToken{"token2", "token3"},
					Commits:      map[graveler.CommitID]*graveler.Commit{"": {}},
				},
				nil, testutil.NewProtectedBranchesManagerFake()),
			expectedValueResult: graveler.Value{Identity: []byte("stagedA")},
		},
		{
			name: "branch - sealed tombstone",
			r: newGraveler(t, &testutil.CommittedFake{Err: graveler.ErrNotFound},
				&testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token2": {"staged": nil},
						"token3": {"staged": {
							Identity: []byte("stagedB"),
							Data:     nil,
						}},
					},
				},
				&testutil.RefsFake{
					RefType:      graveler.ReferenceTypeBranch,
					StagingToken: "token1",
					SealedTokens: []graveler.StagingToken{"token2", "token3"},
					Commits:      map[graveler.CommitID]*graveler.Commit{"": {}},
				},
				nil, testutil.NewProtectedBranchesManagerFake()),
			expectedErr: graveler.ErrNotFound,
		},
		{
			name: "branch -committed, staged entry + tombstone",
			r: newGraveler(t, &testutil.CommittedFake{
				ValuesByKey: map[string]*graveler.Value{"staged": {Identity: []byte("stagedA")}},
			},
				&testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token2": {"staged": nil},
						"token3": {"staged": {
							Identity: []byte("stagedB"),
							Data:     nil,
						}},
					},
				},
				&testutil.RefsFake{
					RefType:      graveler.ReferenceTypeBranch,
					StagingToken: "token1",
					SealedTokens: []graveler.StagingToken{"token2", "token3"},
					Commits:      map[graveler.CommitID]*graveler.Commit{"": {}},
				},
				nil, testutil.NewProtectedBranchesManagerFake()),
			expectedErr: graveler.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Value, err := tt.r.Get(context.Background(), repository, "", []byte("staged"))
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("wrong error, expected:%v got:%v", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expected error
			}
			if string(tt.expectedValueResult.Identity) != string(Value.Identity) {
				t.Errorf("wrong Value address, expected:%s got:%s", tt.expectedValueResult.Identity, Value.Identity)
			}
		})
	}
}

func TestGraveler_Diff(t *testing.T) {
	tests := []struct {
		name            string
		r               catalog.Store
		expectedErr     error
		expectedHasMore bool
		expectedDiff    graveler.DiffIterator
	}{
		{
			name: "no changes",
			r: newGraveler(t, &testutil.CommittedFake{
				DiffIterator: testutil.NewDiffIter([]graveler.Diff{}),
			},
				&testutil.StagingFake{},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID:     "c1",
									StagingToken: "token",
								},
							},
						},
						"ref1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{}),
		},
		{
			name: "no changes - branch staging remove - add",
			r: newGraveler(t, &testutil.CommittedFake{
				Values: map[string]graveler.ValueIterator{
					"mri1": testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}}),
					"mri2": testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{
						Identity: []byte("BAD"),
						Data:     nil,
					}}}),
				},
				DiffIterator: testutil.NewDiffIter([]graveler.Diff{}),
			},
				&testutil.StagingFake{Values: map[string]map[string]*graveler.Value{
					"token": {
						"foo/one": &graveler.Value{},
					},
					"token1": {
						"foo/one": nil,
					},
					"token2": {
						"foo/one": &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						},
					},
				}},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c2", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token1", "token2"}},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}, "c2": {MetaRangeID: "mri2"}},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: graveler.ResolvedBranchModifierStaging,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID:     "c2",
									StagingToken: "token",
									SealedTokens: []graveler.StagingToken{"token1", "token2"},
								},
							},
						},
						"ref1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{}),
		},
		{
			name: "with changes. modify, add, delete",
			r: newGraveler(t, &testutil.CommittedFake{
				Values: map[string]graveler.ValueIterator{
					"mri1": testutil.NewValueIteratorFake([]graveler.ValueRecord{
						{Key: graveler.Key("foo/delete"), Value: &graveler.Value{
							Identity: []byte("deleted"),
							Data:     []byte("deleted"),
						}},
						{Key: graveler.Key("foo/modified_committed"), Value: &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						}},
						{Key: graveler.Key("foo/modify"), Value: &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						}},
					}),
					"mri2": testutil.NewValueIteratorFake([]graveler.ValueRecord{
						{Key: graveler.Key("foo/delete"), Value: &graveler.Value{
							Identity: []byte("deleted"),
							Data:     []byte("deleted"),
						}},
						{Key: graveler.Key("foo/modified_committed"), Value: &graveler.Value{
							Identity: []byte("committed"),
							Data:     []byte("committed"),
						}},
						{Key: graveler.Key("foo/modify"), Value: &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						}},
					}),
				},
				DiffIterator: testutil.NewDiffIter([]graveler.Diff{
					{
						Key:  graveler.Key("foo/modified_committed"),
						Type: graveler.DiffTypeChanged,
						Value: &graveler.Value{
							Identity: []byte("committed"),
							Data:     []byte("committed"),
						},
						LeftIdentity: []byte("DECAF"),
					},
				}),
			},
				&testutil.StagingFake{Values: map[string]map[string]*graveler.Value{
					"token": {
						"foo/add": &graveler.Value{},
					},
					"token1": {
						"foo/delete": nil,
					},
					"token2": {
						"foo/modify": &graveler.Value{
							Identity: []byte("test"),
							Data:     []byte("test"),
						},
					},
				}},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c2", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token1", "token2"}},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}, "c2": {MetaRangeID: "mri2"}},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: graveler.ResolvedBranchModifierStaging,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID:     "c2",
									StagingToken: "token",
									SealedTokens: []graveler.StagingToken{"token1", "token2"},
								},
							},
						},
						"ref1": {
							Type: graveler.ReferenceTypeCommit,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{
				{
					Key:   graveler.Key("foo/add"),
					Type:  graveler.DiffTypeAdded,
					Value: &graveler.Value{},
				},
				{
					Key:  graveler.Key("foo/delete"),
					Type: graveler.DiffTypeRemoved,
					Value: &graveler.Value{
						Identity: []byte("deleted"),
						Data:     []byte("deleted"),
					},
					LeftIdentity: []byte("deleted"),
				},
				{
					Key:  graveler.Key("foo/modified_committed"),
					Type: graveler.DiffTypeChanged,
					Value: &graveler.Value{
						Identity: []byte("committed"),
						Data:     []byte("committed"),
					},
					LeftIdentity: []byte("DECAF"),
				},
				{
					Key:  graveler.Key("foo/modify"),
					Type: graveler.DiffTypeChanged,
					Value: &graveler.Value{
						Identity: []byte("test"),
						Data:     []byte("test"),
					},
					LeftIdentity: []byte("DECAF"),
				},
			}),
		},
		{
			name: "with changes compacted",
			r: newGraveler(t, &testutil.CommittedFake{
				Values: map[string]graveler.ValueIterator{
					"mri1": testutil.NewValueIteratorFake([]graveler.ValueRecord{
						{Key: graveler.Key("foo/delete"), Value: &graveler.Value{
							Identity: []byte("deleted"),
							Data:     []byte("deleted"),
						}},
						{Key: graveler.Key("foo/modified_committed"), Value: &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						}},
						{Key: graveler.Key("foo/modify"), Value: &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						}},
					}),
					"mri2": testutil.NewValueIteratorFake([]graveler.ValueRecord{
						{Key: graveler.Key("foo/delete"), Value: &graveler.Value{
							Identity: []byte("deleted"),
							Data:     []byte("deleted"),
						}},
						{Key: graveler.Key("foo/modified_committed"), Value: &graveler.Value{
							Identity: []byte("committed"),
							Data:     []byte("committed"),
						}},
						{Key: graveler.Key("foo/modify"), Value: &graveler.Value{
							Identity: []byte("DECAF"),
							Data:     []byte("BAD"),
						}},
					}),
				},
				DiffIterator: testutil.NewDiffIter([]graveler.Diff{
					{
						Key:  graveler.Key("foo/modified_committed"),
						Type: graveler.DiffTypeChanged,
						Value: &graveler.Value{
							Identity: []byte("committed"),
							Data:     []byte("committed"),
						},
						LeftIdentity: []byte("DECAF"),
					},
				}),
			},
				&testutil.StagingFake{Values: map[string]map[string]*graveler.Value{
					"token": {
						"foo/add": &graveler.Value{},
					},
					"token1": {
						"foo/delete": nil,
					},
					"token2": {
						"foo/modify": &graveler.Value{
							Identity: []byte("test"),
							Data:     []byte("test"),
						},
					},
				}},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token1", "token2"}, CompactedBaseMetaRangeID: "mri2"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: graveler.ResolvedBranchModifierStaging,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID:     "c1",
									StagingToken: "token",
									SealedTokens: []graveler.StagingToken{"token1", "token2"},
								},
							},
						},
						"ref1": {
							Type: graveler.ReferenceTypeCommit,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{
				{
					Key:   graveler.Key("foo/add"),
					Type:  graveler.DiffTypeAdded,
					Value: &graveler.Value{},
				},
				{
					Key:  graveler.Key("foo/delete"),
					Type: graveler.DiffTypeRemoved,
					Value: &graveler.Value{
						Identity: []byte("deleted"),
						Data:     []byte("deleted"),
					},
					LeftIdentity: []byte("deleted"),
				},
				{
					Key:  graveler.Key("foo/modified_committed"),
					Type: graveler.DiffTypeChanged,
					Value: &graveler.Value{
						Identity: []byte("committed"),
						Data:     []byte("committed"),
					},
					LeftIdentity: []byte("DECAF"),
				},
				{
					Key:  graveler.Key("foo/modify"),
					Type: graveler.DiffTypeChanged,
					Value: &graveler.Value{
						Identity: []byte("test"),
						Data:     []byte("test"),
					},
					LeftIdentity: []byte("DECAF"),
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			diff, err := tt.r.Diff(ctx, repository, "ref1", "b1")
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("wrong error, expected:%s got:%s", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expectedErr
			}

			// compare iterators
			for diff.Next() {
				v := diff.Value()
				if !tt.expectedDiff.Next() {
					t.Fatalf("listing next returned true where expected listing next returned false")
				}
				vEx := tt.expectedDiff.Value()
				require.Nil(t, deep.Equal(v, vEx))
			}
			if tt.expectedDiff.Next() {
				t.Fatalf("expected listing next returned true where listing next returned false")
			}
		})
	}
}

func TestGraveler_DiffUncommitted(t *testing.T) {
	tests := []struct {
		name            string
		r               catalog.Store
		expectedErr     error
		expectedHasMore bool
		expectedDiff    graveler.DiffIterator
	}{
		{
			name: "no changes",
			r: newGraveler(t, &testutil.CommittedFake{
				ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{
					{
						Key: graveler.Key("foo/one"), Value: &graveler.Value{},
					},
				}),
			},
				&testutil.StagingFake{
					ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{}),
				},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{}),
		},
		{
			name: "added one",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{})},
				&testutil.StagingFake{
					ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}}),
				},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{{
				Key:   graveler.Key("foo/one"),
				Type:  graveler.DiffTypeAdded,
				Value: &graveler.Value{},
			}}),
		},
		{
			name: "changed one",
			r: newGraveler(t, &testutil.CommittedFake{
				ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{
					{
						Key: graveler.Key("foo/one"), Value: &graveler.Value{Identity: []byte("one")},
					},
				}),
				ValuesByKey: map[string]*graveler.Value{"foo/one": {Identity: []byte("one")}},
			},
				&testutil.StagingFake{
					ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{Identity: []byte("one_changed")}}}),
					Values: map[string]map[string]*graveler.Value{
						"token": {
							"foo/one": &graveler.Value{Identity: []byte("one_changed")},
						},
					},
				},
				&testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}},
				}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{{
				Key:   graveler.Key("foo/one"),
				Type:  graveler.DiffTypeChanged,
				Value: &graveler.Value{Identity: []byte("one_changed")},
			}}),
		},
		{
			name: "removed one",
			r: newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{Identity: []byte("not-nil")}}})},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: nil}, {Key: graveler.Key("foo/two"), Value: nil}})},
				&testutil.RefsFake{Branch: &graveler.Branch{CommitID: "c1", StagingToken: "token"}, Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{{
				Key:  graveler.Key("foo/one"),
				Type: graveler.DiffTypeRemoved,
			}}),
		},
		{
			name: "diff with compacted",
			r: newGraveler(t,
				&testutil.CommittedFake{
					Values: map[string]graveler.ValueIterator{
						"mri1": testutil.NewValueIteratorFake([]graveler.ValueRecord{
							{Key: graveler.Key("foo/a"), Value: &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BAD"),
							}},
							{Key: graveler.Key("foo/b"), Value: &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BAD"),
							}},
							{Key: graveler.Key("foo/c"), Value: &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BAD"),
							}},
							{Key: graveler.Key("foo/d"), Value: &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BAD"),
							}},
							{Key: graveler.Key("foo/e"), Value: &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BAD"),
							}},
							{Key: graveler.Key("foo/g"), Value: &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BAD"),
							}},
						}),
					},
					DiffIterator: testutil.NewDiffIter([]graveler.Diff{
						{
							Key:   graveler.Key("foo/a"),
							Type:  graveler.DiffTypeRemoved,
							Value: &graveler.Value{},
						},
						{
							Key:   graveler.Key("foo/b"),
							Type:  graveler.DiffTypeChanged,
							Value: &graveler.Value{Identity: []byte("compacted"), Data: []byte("compacted")},
						},
						{
							Key:  graveler.Key("foo/c"),
							Type: graveler.DiffTypeRemoved,
						},
						{
							Key:  graveler.Key("foo/d"),
							Type: graveler.DiffTypeRemoved,
						},
						{
							Key:   graveler.Key("foo/e"),
							Type:  graveler.DiffTypeAdded,
							Value: &graveler.Value{Identity: []byte("BAD"), Data: []byte("BAD")},
						},
						{
							Key:   graveler.Key("foo/f"),
							Type:  graveler.DiffTypeChanged,
							Value: &graveler.Value{Identity: []byte("BAD"), Data: []byte("BAD")},
						},
						{
							Key:   graveler.Key("foo/g"),
							Type:  graveler.DiffTypeChanged,
							Value: &graveler.Value{Identity: []byte("BAD"), Data: []byte("BAD")},
						},
					}),
				},
				&testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{
					{Key: graveler.Key("foo/d"), Value: &graveler.Value{Identity: []byte("staged"), Data: []byte("staged")}},
					{Key: graveler.Key("foo/e"), Value: &graveler.Value{Identity: []byte("staged"), Data: []byte("staged")}},
					{Key: graveler.Key("foo/f"), Value: &graveler.Value{Identity: []byte("staged"), Data: []byte("staged")}},
					{Key: graveler.Key("foo/g"), Value: nil},
				})},
				&testutil.RefsFake{Branch: &graveler.Branch{CommitID: "c1", StagingToken: "token", CompactedBaseMetaRangeID: "mri2"}, Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mri1"}}}, nil, testutil.NewProtectedBranchesManagerFake(),
			),
			expectedDiff: testutil.NewDiffIter([]graveler.Diff{
				{
					Key:   graveler.Key("foo/a"),
					Type:  graveler.DiffTypeRemoved,
					Value: &graveler.Value{},
				},
				{
					Key:   graveler.Key("foo/b"),
					Type:  graveler.DiffTypeChanged,
					Value: &graveler.Value{Identity: []byte("compacted"), Data: []byte("compacted")},
				},
				{
					Key:  graveler.Key("foo/c"),
					Type: graveler.DiffTypeRemoved,
				},
				{
					Key:   graveler.Key("foo/d"),
					Type:  graveler.DiffTypeChanged,
					Value: &graveler.Value{Identity: []byte("staged"), Data: []byte("staged")},
				},
				{
					Key:   graveler.Key("foo/e"),
					Type:  graveler.DiffTypeChanged,
					Value: &graveler.Value{Identity: []byte("staged"), Data: []byte("staged")},
				},
				{
					Key:   graveler.Key("foo/f"),
					Type:  graveler.DiffTypeAdded,
					Value: &graveler.Value{Identity: []byte("staged"), Data: []byte("staged")},
				},
				{
					Key:  graveler.Key("foo/g"),
					Type: graveler.DiffTypeRemoved,
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			diff, err := tt.r.DiffUncommitted(ctx, repository, "branch")
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("wrong error, expected:%s got:%s", tt.expectedErr, err)
			}
			if err != nil {
				return // err == tt.expectedErr
			}

			// compare iterators
			for diff.Next() {
				if !tt.expectedDiff.Next() {
					t.Fatalf("listing next returned true where expected listing next returned false")
				}
				if diff := deep.Equal(diff.Value(), tt.expectedDiff.Value()); diff != nil {
					t.Errorf("unexpected diff %s", diff)
				}
			}
			if tt.expectedDiff.Next() {
				t.Fatalf("expected listing next returned true where listing next returned false")
			}
		})
	}
}

func TestGravelerDiffUncommitted_Advanced(t *testing.T) {
	committedFake := &testutil.CommittedFake{
		ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{
			{
				Key: graveler.Key("in_staged_actual_no_change"), Value: &graveler.Value{Identity: []byte("test1")},
			},
			{
				Key: graveler.Key("not_in_staged_in_sealed"), Value: &graveler.Value{Identity: []byte("BAD")},
			},
			{
				Key: graveler.Key("staged"), Value: &graveler.Value{Identity: []byte("BAD")},
			},
		}),
	}
	stagingFake := &testutil.StagingFake{
		Values: map[string]map[string]*graveler.Value{
			"token": {
				"in_staged_actual_no_change": {
					Identity: []byte("test1"),
					Data:     nil,
				},
				"staged": {
					Identity: []byte("stagedA"),
					Data:     nil,
				},
			},
			"token1": {
				"added_in_staging": {
					Identity: []byte("stagedC"),
					Data:     nil,
				},
				"in_staged_actual_no_change": {
					Identity: []byte("DEAD"),
					Data:     nil,
				},
				"not_in_staged_in_sealed": {
					Identity: []byte("stagedB"),
					Data:     nil,
				},
				"staged": {
					Identity: []byte("DEAD"),
					Data:     nil,
				},
			},
			"token2": {
				"in_staged_actual_no_change": {
					Identity: []byte("BEEF"),
					Data:     nil,
				},
				"not_in_staged_in_sealed": {
					Identity: []byte("DECAF"),
					Data:     nil,
				},
				"staged": {
					Identity: []byte("BEEF"),
					Data:     nil,
				},
			},
		},
	}
	refsFake := &testutil.RefsFake{
		RefType:      graveler.ReferenceTypeBranch,
		StagingToken: "token",
		Branch:       &graveler.Branch{CommitID: "c1", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token1", "token2"}},
		Commits:      map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mr1"}},
	}
	r := newGraveler(t, committedFake,
		stagingFake,
		refsFake,
		nil, testutil.NewProtectedBranchesManagerFake())
	expectedDiff := testutil.NewDiffIter([]graveler.Diff{
		{
			Key:  graveler.Key("added_in_staging"),
			Type: graveler.DiffTypeAdded,
			Value: &graveler.Value{
				Identity: []byte("stagedC"),
			},
		},
		{
			Key:  graveler.Key("not_in_staged_in_sealed"),
			Type: graveler.DiffTypeChanged,
			Value: &graveler.Value{
				Identity: []byte("stagedB"),
			},
		},
		{
			Key:  graveler.Key("staged"),
			Type: graveler.DiffTypeChanged,
			Value: &graveler.Value{
				Identity: []byte("stagedA"),
			},
		},
	})

	ctx := context.Background()
	diff, err := r.DiffUncommitted(ctx, repository, "branch")
	require.NoError(t, err)
	// compare iterators
	for diff.Next() {
		v := diff.Value()
		if !expectedDiff.Next() {
			t.Fatalf("listing next returned true where expected listing next returned false")
		}
		exV := expectedDiff.Value()
		if deep.Equal(v, exV) != nil {
			t.Errorf("unexpected diff actual: %v, expected: %v", v, exV)
		}
	}
	if expectedDiff.Next() {
		t.Fatalf("expected listing next returned true where listing next returned false")
	}
}

func TestGraveler_CreateBranch(t *testing.T) {
	gravel := newGraveler(t, nil, nil, &testutil.RefsFake{Err: graveler.ErrBranchNotFound, CommitID: "8888888798e3aeface8e62d1c7072a965314b4"}, nil, nil)
	_, err := gravel.CreateBranch(context.Background(), repository, "", "")
	if err != nil {
		t.Fatal("unexpected error on create branch", err)
	}
	// test create branch when branch exists
	gravel = newGraveler(t, nil, nil, &testutil.RefsFake{Branch: &graveler.Branch{}}, nil, nil)
	_, err = gravel.CreateBranch(context.Background(), repository, "", "")
	if !errors.Is(err, graveler.ErrBranchExists) {
		t.Fatal("did not get expected error, expected ErrBranchExists")
	}
}

func TestGraveler_UpdateBranch(t *testing.T) {
	gravel := newGraveler(t, nil, &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: graveler.Key("foo/one"), Value: &graveler.Value{}}})},
		&testutil.RefsFake{Branch: &graveler.Branch{}, UpdateErr: kv.ErrPredicateFailed}, nil, nil)
	testutil.ShortenBranchUpdateBackOff(gravel.(*graveler.Graveler))
	_, err := gravel.UpdateBranch(context.Background(), repository, "", "")
	require.ErrorIs(t, err, graveler.ErrTooManyTries)

	gravel = newGraveler(t, &testutil.CommittedFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{})}, &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{})},
		&testutil.RefsFake{Branch: &graveler.Branch{StagingToken: "st1", CommitID: "commit1"}, Commits: map[graveler.CommitID]*graveler.Commit{"commit1": {}}}, nil, nil)
	_, err = gravel.UpdateBranch(context.Background(), repository, "", "")
	require.NoError(t, err)
}

func TestGravelerCommit(t *testing.T) {
	expectedCommitID := graveler.CommitID("expectedCommitId")
	expectedRangeID := graveler.MetaRangeID("expectedRangeID")
	values := testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: nil, Value: nil}})
	multipleValues := []graveler.ValueIterator{
		testutil.NewValueIteratorFake([]graveler.ValueRecord{}),
		testutil.NewValueIteratorFake([]graveler.ValueRecord{}),
	}
	type fields struct {
		CommittedManager         *testutil.CommittedFake
		StagingManager           *testutil.StagingFake
		RefManager               *testutil.RefsFake
		ProtectedBranchesManager *testutil.ProtectedBranchesManagerFake
	}
	type args struct {
		ctx             context.Context
		branchID        graveler.BranchID
		committer       string
		message         string
		metadata        graveler.Metadata
		sourceMetarange *graveler.MetaRangeID
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        graveler.CommitID
		values      graveler.ValueIterator
		expectedErr error
	}{
		{
			name: "valid commit without source metarange",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				branchID:  "branch",
				committer: "committer",
				message:   "a message",
				metadata:  graveler.Metadata{},
			},
			want:        expectedCommitID,
			values:      values,
			expectedErr: nil,
		},
		{
			name: "valid commit with source metarange",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager: &testutil.StagingFake{
					ValueIterator: testutil.NewValueIteratorFake([]graveler.ValueRecord{}),
				},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID, StagingToken: "token1"},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:             nil,
				branchID:        "branch",
				committer:       "committer",
				message:         "a message",
				metadata:        graveler.Metadata{},
				sourceMetarange: &expectedRangeID,
			},
			want:        expectedCommitID,
			values:      values,
			expectedErr: nil,
		},
		{
			name: "commit with source metarange and non-empty staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager: &testutil.StagingFake{ValueIterator: testutils.NewFakeValueIterator([]*graveler.ValueRecord{{
					Key: []byte("key1"), Value: &graveler.Value{Identity: []byte("id1"), Data: []byte("data1")},
				}})},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID, StagingToken: "token1"},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:             nil,
				branchID:        "branch",
				committer:       "committer",
				message:         "a message",
				metadata:        graveler.Metadata{},
				sourceMetarange: &expectedRangeID,
			},
			values:      values,
			expectedErr: graveler.ErrCommitMetaRangeDirtyBranch,
		},
		{
			name: "commit with source metarange and non-empty compaction",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID, DiffIterator: testutil.NewDiffIter([]graveler.Diff{{Key: key1, Type: graveler.DiffTypeRemoved}})},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutils.NewFakeValueIterator([]*graveler.ValueRecord{})},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID, StagingToken: "token1", CompactedBaseMetaRangeID: mr2ID},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:             nil,
				branchID:        "branch",
				committer:       "committer",
				message:         "a message",
				metadata:        graveler.Metadata{},
				sourceMetarange: &expectedRangeID,
			},
			values:      values,
			expectedErr: graveler.ErrCommitMetaRangeDirtyBranch,
		},
		{
			name: "fail on apply",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID, Err: graveler.ErrConflictFound},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				branchID:  "branch",
				committer: "committer",
				message:   "a message",
				metadata:  nil,
			},
			want:        expectedCommitID,
			values:      values,
			expectedErr: graveler.ErrConflictFound,
		},
		{
			name: "fail on add commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{
					CommitID:  expectedCommitID,
					Branch:    &graveler.Branch{CommitID: expectedCommitID},
					CommitErr: graveler.ErrConflictFound,
					Commits:   map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				branchID:  "branch",
				committer: "committer",
				message:   "a message",
				metadata:  nil,
			},
			want:        expectedCommitID,
			values:      values,
			expectedErr: graveler.ErrConflictFound,
		},
		{
			name: "fail on drop",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values, DropErr: graveler.ErrNotFound},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				branchID:  "branch",
				committer: "committer",
				message:   "a message",
				metadata:  graveler.Metadata{},
			},
			want:        expectedCommitID,
			values:      values,
			expectedErr: nil,
		},
		{
			name: "fail on protected branch",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
				ProtectedBranchesManager: testutil.NewProtectedBranchesManagerFake("branch"),
			},
			args: args{
				ctx:       nil,
				branchID:  "branch",
				committer: "committer",
				message:   "a message",
				metadata:  graveler.Metadata{},
			},
			values:      values,
			expectedErr: graveler.ErrCommitToProtectedBranch,
		},
		{
			name: "valid commit with staging and sealed",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: values},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID,
					Branch:   &graveler.Branch{CommitID: expectedCommitID, StagingToken: "token", SealedTokens: []graveler.StagingToken{"token", "token2"}},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				branchID:  "branch",
				committer: "committer",
				message:   "a message",
				metadata:  graveler.Metadata{},
			},
			want:        expectedCommitID,
			values:      graveler.NewCombinedIterator(multipleValues...),
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedCommitID = "expectedCommitId"
			expectedRangeID = "expectedRangeID"
			if tt.fields.ProtectedBranchesManager == nil {
				tt.fields.ProtectedBranchesManager = testutil.NewProtectedBranchesManagerFake()
			}
			g := newGraveler(t, tt.fields.CommittedManager, tt.fields.StagingManager, tt.fields.RefManager, nil, tt.fields.ProtectedBranchesManager)

			got, err := g.Commit(context.Background(), repository, tt.args.branchID, graveler.CommitParams{
				Committer:       tt.args.committer,
				Message:         tt.args.message,
				Metadata:        tt.args.metadata,
				SourceMetaRange: tt.args.sourceMetarange,
			})
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("unexpected err got = %v, wanted = %v", err, tt.expectedErr)
			}
			if err != nil {
				return
			}
			expectedAppliedData := testutil.AppliedData{
				Values:      tt.values,
				MetaRangeID: expectedRangeID,
			}
			if tt.args.sourceMetarange != nil {
				expectedAppliedData = testutil.AppliedData{}
			}
			if diff := deep.Equal(tt.fields.CommittedManager.AppliedData, expectedAppliedData); diff != nil {
				t.Errorf("unexpected apply data %s", diff)
			}

			if diff := deep.Equal(tt.fields.RefManager.AddedCommit, testutil.AddedCommitData{
				Committer:   tt.args.committer,
				Message:     tt.args.message,
				MetaRangeID: expectedRangeID,
				Parents:     graveler.CommitParents{expectedCommitID},
				Metadata:    graveler.Metadata{},
			}); diff != nil {
				t.Errorf("unexpected added commit %s", diff)
			}
			if !tt.fields.StagingManager.DropCalled && tt.fields.StagingManager.DropErr == nil {
				t.Errorf("expected drop to be called")
			}

			if got != expectedCommitID {
				t.Errorf("got wrong commitID, got = %v, want %v", got, expectedCommitID)
			}
		})
	}
}

// TestGraveler_MergeInvalidRef test merge with invalid source reference in order
func TestGraveler_MergeInvalidRef(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const destinationCommitID = graveler.CommitID("destinationCommitID")
	const mergeDestination = graveler.BranchID("destinationID")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		Err:    graveler.ErrInvalidRef,
		Branch: &graveler.Branch{CommitID: destinationCommitID, StagingToken: "st1"},
		Refs: map[graveler.Ref]*graveler.ResolvedRef{
			graveler.Ref(mergeDestination): {
				Type: graveler.ReferenceTypeBranch,
				BranchRecord: graveler.BranchRecord{
					BranchID: mergeDestination,
					Branch: &graveler.Branch{
						CommitID: destinationCommitID,
					},
				},
			},
		},
		Commits: map[graveler.CommitID]*graveler.Commit{
			destinationCommitID: {MetaRangeID: expectedRangeID},
		},
	}
	g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())

	// test merge invalid ref
	ctx := context.Background()
	const commitCommitter = "committer"
	const mergeMessage = "message"
	_, err := g.Merge(ctx, repository, mergeDestination, "unexpectedRef", graveler.CommitParams{
		Committer: commitCommitter,
		Message:   mergeMessage,
		Metadata:  graveler.Metadata{"key1": "val1"},
	}, "")
	if !errors.Is(err, graveler.ErrInvalidRef) {
		t.Fatalf("Merge failed with err=%v, expected ErrInvalidRef", err)
	}
}

func TestGraveler_AddCommit(t *testing.T) {
	const (
		expectedCommitID       = graveler.CommitID("expectedCommitId")
		expectedParentCommitID = graveler.CommitID("expectedParentCommitId")
		expectedRangeID        = graveler.MetaRangeID("expectedRangeID")
	)

	type fields struct {
		CommittedManager *testutil.CommittedFake
		StagingManager   *testutil.StagingFake
		RefManager       *testutil.RefsFake
	}
	type args struct {
		ctx            context.Context
		committer      string
		message        string
		metadata       graveler.Metadata
		missingParents bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        graveler.CommitID
		expectedErr error
	}{
		{
			name: "valid commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)},
				RefManager: &testutil.RefsFake{
					CommitID: expectedCommitID, Branch: &graveler.Branch{CommitID: expectedParentCommitID},
					Commits: map[graveler.CommitID]*graveler.Commit{
						expectedParentCommitID: {},
					},
				},
			},
			args: args{
				ctx:       nil,
				committer: "committer",
				message:   "a message",
				metadata:  graveler.Metadata{},
			},
			want:        expectedCommitID,
			expectedErr: nil,
		},
		{
			name: "meta range not found",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{Err: graveler.ErrMetaRangeNotFound},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)},
				RefManager: &testutil.RefsFake{
					CommitID: expectedParentCommitID,
					Branch:   &graveler.Branch{CommitID: expectedParentCommitID},
					Commits:  map[graveler.CommitID]*graveler.Commit{expectedParentCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				committer: "committer",
				message:   "a message",
				metadata:  nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrMetaRangeNotFound,
		},
		{
			name: "fail on add commit",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{MetaRangeID: expectedRangeID},
				StagingManager:   &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)},
				RefManager: &testutil.RefsFake{
					CommitID:  expectedCommitID,
					Branch:    &graveler.Branch{CommitID: expectedParentCommitID},
					CommitErr: graveler.ErrConflictFound,
					Commits:   map[graveler.CommitID]*graveler.Commit{expectedParentCommitID: {MetaRangeID: expectedRangeID}},
				},
			},
			args: args{
				ctx:       nil,
				committer: "committer",
				message:   "a message",
				metadata:  nil,
			},
			want:        expectedCommitID,
			expectedErr: graveler.ErrConflictFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGraveler(t, tt.fields.CommittedManager, tt.fields.StagingManager, tt.fields.RefManager, nil, testutil.NewProtectedBranchesManagerFake())
			commit := graveler.Commit{
				Committer:   tt.args.committer,
				Message:     tt.args.message,
				MetaRangeID: expectedRangeID,
				Metadata:    tt.args.metadata,
			}
			if !tt.args.missingParents {
				commit.Parents = graveler.CommitParents{expectedParentCommitID}
			}
			got, err := g.AddCommit(context.Background(), repository, commit)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("unexpected err got = %v, wanted = %v", err, tt.expectedErr)
			}
			if err != nil {
				return
			}

			if diff := deep.Equal(tt.fields.RefManager.AddedCommit, testutil.AddedCommitData{
				Committer:   tt.args.committer,
				Message:     tt.args.message,
				MetaRangeID: expectedRangeID,
				Parents:     graveler.CommitParents{expectedParentCommitID},
				Metadata:    graveler.Metadata{},
			}); diff != nil {
				t.Errorf("unexpected added commit %s", diff)
			}
			if tt.fields.StagingManager.DropCalled {
				t.Error("Staging manager drop shouldn't be called")
			}

			if got != expectedCommitID {
				t.Errorf("got wrong commitID, got = %v, want %v", got, expectedCommitID)
			}
		})
	}
}

func TestGravelerDelete(t *testing.T) {
	type fields struct {
		CommittedManager graveler.CommittedManager
		StagingManager   *testutil.StagingFake
		RefManager       graveler.RefManager
	}
	type args struct {
		branchID graveler.BranchID
		key      graveler.Key
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedSetValue   *graveler.ValueRecord
		expectedRemovedKey graveler.Key
		expectedErr        error
	}{
		{
			name: "exists only in committed",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key": {}},
				},
				StagingManager: &testutil.StagingFake{
					Err: graveler.ErrNotFound,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args: args{
				key: []byte("key"),
			},
			expectedSetValue: &graveler.ValueRecord{
				Key:   []byte("key"),
				Value: nil,
			},
			expectedErr: nil,
		},
		{
			name: "exists only in compacted",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key": {}},
				},
				StagingManager: &testutil.StagingFake{
					Err: graveler.ErrNotFound,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", CompactedBaseMetaRangeID: "mr2"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mr1"}},
				},
			},
			args: args{
				key: []byte("key"),
			},
			expectedSetValue: &graveler.ValueRecord{
				Key:   []byte("key"),
				Value: nil,
			},
			expectedErr: nil,
		},
		{
			name: "exists in committed and in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key1": {}},
				},
				StagingManager: &testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token": {
							"key2": &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BEEF"),
							},
						},
						"token2": {
							"key1": &graveler.Value{
								Identity: []byte("test"),
								Data:     []byte("test"),
							},
						},
					},
					Value: &graveler.Value{},
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token", "token2"}},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args: args{
				key: []byte("key1"),
			},
			expectedSetValue: &graveler.ValueRecord{
				Key:   []byte("key1"),
				Value: nil,
			},
			expectedErr: nil,
		},
		{
			name: "exists in compacted and in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key1": {}},
				},
				StagingManager: &testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token": {
							"key2": &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BEEF"),
							},
						},
						"token2": {
							"key1": &graveler.Value{
								Identity: []byte("test"),
								Data:     []byte("test"),
							},
						},
					},
					Value: &graveler.Value{},
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token", "token2"}, CompactedBaseMetaRangeID: "mr2"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mr1"}},
				},
			},
			args: args{
				key: []byte("key1"),
			},
			expectedSetValue: &graveler.ValueRecord{
				Key:   []byte("key1"),
				Value: nil,
			},
			expectedErr: nil,
		},
		{
			name: "exists in committed tombstone in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key1": {}},
				},
				StagingManager: &testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token": {
							"key1": nil,
						},
						"token2": {
							"key1": &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BEEF"),
							},
						},
					},
					Value: nil,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token", "token2"}},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args:        args{key: []byte("key1")},
			expectedErr: nil,
		},
		{
			name: "exists in compacted tombstone in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					ValuesByKey: map[string]*graveler.Value{"key1": {}},
				},
				StagingManager: &testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{
						"token": {
							"key1": nil,
						},
						"token2": {
							"key1": &graveler.Value{
								Identity: []byte("BAD"),
								Data:     []byte("BEEF"),
							},
						},
					},
					Value: nil,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token", SealedTokens: []graveler.StagingToken{"token", "token2"}, CompactedBaseMetaRangeID: "mr2"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mr1"}},
				},
			},
			args:        args{key: []byte("key1")},
			expectedErr: nil,
		},
		{
			name: "exists only in staging - commits",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					Err: graveler.ErrNotFound,
				},
				StagingManager: &testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{"token": {"key1": &graveler.Value{
						Identity: []byte("test"),
						Data:     []byte("test"),
					}}},
					Value: nil,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {}},
				},
			},
			args: args{
				key: []byte("key1"),
			},
			expectedRemovedKey: []byte("key1"),
			expectedErr:        nil,
		},
		{
			name: "exists in staging and not in compaction",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					Err: graveler.ErrNotFound,
				},
				StagingManager: &testutil.StagingFake{
					Values: map[string]map[string]*graveler.Value{"token": {"key1": &graveler.Value{
						Identity: []byte("test"),
						Data:     []byte("test"),
					}}},
					Value: nil,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CommitID: "c1", StagingToken: "token", CompactedBaseMetaRangeID: "mr2"},
					Commits: map[graveler.CommitID]*graveler.Commit{"c1": {MetaRangeID: "mr1"}},
				},
			},
			args: args{
				key: []byte("key1"),
			},
			expectedRemovedKey: []byte("key1"),
			expectedErr:        nil,
		},
		{
			name: "not in committed not in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					Err: graveler.ErrNotFound,
				},
				StagingManager: &testutil.StagingFake{
					Err: graveler.ErrNotFound,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{},
					Commits: map[graveler.CommitID]*graveler.Commit{"": {}},
				},
			},
			args:        args{},
			expectedErr: nil,
		},
		{
			name: "not in compacted not in staging",
			fields: fields{
				CommittedManager: &testutil.CommittedFake{
					Err: graveler.ErrNotFound,
				},
				StagingManager: &testutil.StagingFake{
					Err: graveler.ErrNotFound,
				},
				RefManager: &testutil.RefsFake{
					Branch:  &graveler.Branch{CompactedBaseMetaRangeID: "mr2"},
					Commits: map[graveler.CommitID]*graveler.Commit{"": {MetaRangeID: "mr1"}},
				},
			},
			args:        args{},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			g := newGraveler(t, tt.fields.CommittedManager, tt.fields.StagingManager, tt.fields.RefManager, nil, testutil.NewProtectedBranchesManagerFake())
			if err := g.Delete(ctx, repository, tt.args.branchID, tt.args.key); !errors.Is(err, tt.expectedErr) {
				t.Errorf("Delete() returned unexpected error. got = %v, expected %v", err, tt.expectedErr)
			}

			if tt.expectedRemovedKey != nil {
				// validate set on staging
				if diff := deep.Equal(tt.fields.StagingManager.LastSetValueRecord, &graveler.ValueRecord{Key: tt.expectedRemovedKey, Value: nil}); diff != nil {
					t.Errorf("unexpected set value %s", diff)
				}
			}
		})
	}
}

func TestGraveler_PrepareCommitHook(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const expectedCommitID = graveler.CommitID("expectedCommitId")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		CommitID: expectedCommitID,
		Branch:   &graveler.Branch{CommitID: expectedCommitID},
		Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
	}
	// tests
	errSomethingBad := errors.New("something bad")
	const commitBranchID = "branchID"
	const commitCommitter = "committer"
	const commitMessage = "message"
	commitMetadata := graveler.Metadata{"key1": "val1"}
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name:         "without hook",
			hook:         false,
			err:          nil,
			readOnlyRepo: false,
		},
		{
			name:         "hook no error",
			hook:         true,
			err:          nil,
			readOnlyRepo: false,
		},
		{
			name:         "hook read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
		{
			name:         "hook error",
			hook:         true,
			err:          errSomethingBad,
			readOnlyRepo: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreCommitHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			// call commit
			_, err := g.Commit(ctx, repo, commitBranchID, graveler.CommitParams{
				Committer: commitCommitter,
				Message:   commitMessage,
				Metadata:  commitMetadata,
			}, graveler.WithForce(tt.readOnlyRepo))

			// check err composition
			if !errors.Is(err, tt.err) {
				t.Fatalf("Commit err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Commit err=%v, expected HookAbortError", err)
			}
			called := slices.Contains(h.Called, "PrepareCommitHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Commit invalid prepare-commit hook call, %v expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repo.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repo.RepositoryID)
			}
			if h.BranchID != commitBranchID {
				t.Errorf("Hook branch '%s', expected '%s'", h.BranchID, commitBranchID)
			}
			if h.Commit.Message != commitMessage {
				t.Errorf("Hook commit message '%s', expected '%s'", h.Commit.Message, commitMessage)
			}
			if diff := deep.Equal(h.Commit.Metadata, commitMetadata); diff != nil {
				t.Error("Hook commit metadata diff:", diff)
			}
		})
	}
}

func TestGraveler_PreCommitHook(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const expectedCommitID = graveler.CommitID("expectedCommitId")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		CommitID: expectedCommitID,
		Branch:   &graveler.Branch{CommitID: expectedCommitID},
		Commits:  map[graveler.CommitID]*graveler.Commit{expectedCommitID: {MetaRangeID: expectedRangeID}},
	}
	// tests
	errSomethingBad := errors.New("something bad")
	const commitBranchID = "branchID"
	const commitCommitter = "committer"
	const commitMessage = "message"
	commitMetadata := graveler.Metadata{"key1": "val1"}
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name:         "without hook",
			hook:         false,
			err:          nil,
			readOnlyRepo: false,
		},
		{
			name:         "hook no error",
			hook:         true,
			err:          nil,
			readOnlyRepo: false,
		},
		{
			name:         "hook read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
		{
			name:         "hook error",
			hook:         true,
			err:          errSomethingBad,
			readOnlyRepo: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreCommitHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			// call commit
			_, err := g.Commit(ctx, repo, commitBranchID, graveler.CommitParams{
				Committer: commitCommitter,
				Message:   commitMessage,
				Metadata:  commitMetadata,
			}, graveler.WithForce(tt.readOnlyRepo))

			// check err composition
			if !errors.Is(err, tt.err) {
				t.Fatalf("Commit err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Commit err=%v, expected HookAbortError", err)
			}
			called := slices.Contains(h.Called, "PrepareCommitHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Commit invalid pre-hook call, %v expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repo.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repo.RepositoryID)
			}
			if h.BranchID != commitBranchID {
				t.Errorf("Hook branch '%s', expected '%s'", h.BranchID, commitBranchID)
			}
			if h.Commit.Message != commitMessage {
				t.Errorf("Hook commit message '%s', expected '%s'", h.Commit.Message, commitMessage)
			}
			if diff := deep.Equal(h.Commit.Metadata, commitMetadata); diff != nil {
				t.Error("Hook commit metadata diff:", diff)
			}
		})
	}
}

func TestGraveler_PreMergeHook(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const expectedCommitID = graveler.CommitID("expectedCommitID")
	const destinationCommitID = graveler.CommitID("destinationCommitID")
	const mergeDestination = graveler.BranchID("destinationID")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		CommitID: expectedCommitID,
		Branch:   &graveler.Branch{CommitID: destinationCommitID, StagingToken: "st1"},
		Refs: map[graveler.Ref]*graveler.ResolvedRef{
			graveler.Ref(mergeDestination): {
				Type: graveler.ReferenceTypeBranch,
				BranchRecord: graveler.BranchRecord{
					BranchID: mergeDestination,
					Branch: &graveler.Branch{
						CommitID:     destinationCommitID,
						StagingToken: "st2",
					},
				},
			},
		},
		Commits: map[graveler.CommitID]*graveler.Commit{
			expectedCommitID:    {MetaRangeID: expectedRangeID},
			destinationCommitID: {MetaRangeID: expectedRangeID},
		},
	}
	// tests
	errSomethingBad := errors.New("first error")
	const commitCommitter = "committer"
	const mergeMessage = "message"
	mergeMetadata := graveler.Metadata{"key1": "val1"}
	expectedMergeMetadata := graveler.Metadata{
		"key1":                   "val1",
		".lakefs.merge.strategy": "default",
	}
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name:         "without hook",
			hook:         false,
			err:          nil,
			readOnlyRepo: false,
		},
		{
			name:         "hook no error",
			hook:         true,
			err:          nil,
			readOnlyRepo: false,
		},
		{
			name:         "hook read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
		{
			name:         "hook error",
			hook:         true,
			err:          errSomethingBad,
			readOnlyRepo: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreMergeHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			// call merge
			mergeCommitID, err := g.Merge(ctx, repo, mergeDestination, expectedCommitID.Ref(), graveler.CommitParams{
				Committer: commitCommitter,
				Message:   mergeMessage,
				Metadata:  mergeMetadata,
			}, "", graveler.WithForce(tt.readOnlyRepo))
			// verify we got an error
			if !errors.Is(err, tt.err) {
				t.Fatalf("Merge err=%v, pre-merge error expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Merge err=%v, pre-merge error expected HookAbortError", err)
			}
			if refManager.AddedCommit.MetaRangeID == "" {
				t.Fatalf("Empty MetaRangeID, commit was successful - %+v", refManager.AddedCommit)
			}
			parents := refManager.AddedCommit.Parents
			if len(parents) != 2 {
				t.Fatalf("Merge commit should have 2 parents (%v)", parents)
			}
			if parents[0] != destinationCommitID || parents[1] != expectedCommitID {
				t.Fatalf("Wrong CommitParents order, expected: (%s, %s), got: (%s, %s)", destinationCommitID, expectedCommitID, parents[0], parents[1])
			}
			// verify that calls made until the first error
			called := slices.Contains(h.Called, "PreMergeHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Merge hook h.Called=%v, expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repository.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repository.RepositoryID)
			}
			if h.BranchID != mergeDestination {
				t.Errorf("Hook branch (destination) '%s', expected '%s'", h.BranchID, mergeDestination)
			}
			if h.SourceRef.String() != expectedCommitID.String() {
				t.Errorf("Hook source '%s', expected '%s'", h.SourceRef, expectedCommitID)
			}
			if h.Commit.Message != mergeMessage {
				t.Errorf("Hook merge message '%s', expected '%s'", h.Commit.Message, mergeMessage)
			}
			if h.CommitID != mergeCommitID {
				t.Errorf("Hook merge commit ID '%s', expected '%s'", h.CommitID, mergeCommitID)
			}
			if diff := deep.Equal(h.Commit.Metadata, expectedMergeMetadata); diff != nil {
				t.Error("Hook merge metadata diff:", diff)
			}
		})
	}
}

func TestGraveler_CreateTag(t *testing.T) {
	// prepare graveler
	const commitID = graveler.CommitID("commitID")
	const tagID = graveler.TagID("tagID")
	committedManager := &testutil.CommittedFake{}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		Err: graveler.ErrTagNotFound,
	}
	// tests
	errSomethingBad := errors.New("first error")
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Successful",
			err:  nil,
		},
		{
			name: "Tag exists",
			err:  graveler.ErrTagAlreadyExists,
		},
		{
			name: "Other error on get tag",
			err:  errSomethingBad,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()

			if tt.err != nil {
				refManager.Err = tt.err
			}
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			err := g.CreateTag(ctx, repository, tagID, commitID)

			// verify we got an error
			if !errors.Is(err, tt.err) {
				t.Fatalf("Create tag err=%v, expected=%v", err, tt.err)
			}
		})
	}
}

func TestGraveler_PreCreateTagHook(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const expectedCommitID = graveler.CommitID("expectedCommitID")
	const expectedTagID = graveler.TagID("expectedTagID")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		CommitID: expectedCommitID,
		Branch:   &graveler.Branch{CommitID: expectedCommitID},
		Err:      graveler.ErrTagNotFound,
		Commits: map[graveler.CommitID]*graveler.Commit{
			expectedCommitID: {MetaRangeID: expectedRangeID},
		},
	}
	// tests
	errSomethingBad := errors.New("first error")
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name: "without hook",
			hook: false,
			err:  nil,
		},
		{
			name: "hook no error",
			hook: true,
			err:  nil,
		},
		{
			name: "hook error",
			hook: true,
			err:  errSomethingBad,
		},
		{
			name:         "read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreCreateTagHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}

			err := g.CreateTag(ctx, repo, expectedTagID, expectedCommitID, graveler.WithForce(tt.readOnlyRepo))

			// verify we got an error
			if !errors.Is(err, tt.err) {
				t.Fatalf("Create tag err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Create tag err=%v, expected HookAbortError", err)
			}

			// verify that calls made until the first error
			called := slices.Contains(h.Called, "PreCreateTagHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Pre-create tag hook h.Called=%v, expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repository.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repository.RepositoryID)
			}
			if h.CommitID != expectedCommitID {
				t.Errorf("Hook commit ID '%s', expected '%s'", h.BranchID, expectedCommitID)
			}
			if h.TagID != expectedTagID {
				t.Errorf("Hook tag ID '%s', expected '%s'", h.TagID, expectedTagID)
			}
		})
	}
}

func TestGraveler_PreDeleteTagHook(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const expectedCommitID = graveler.CommitID("expectedCommitID")
	const expectedTagID = graveler.TagID("expectedTagID")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		CommitID: expectedCommitID,
		Branch:   &graveler.Branch{CommitID: expectedCommitID},
		Commits: map[graveler.CommitID]*graveler.Commit{
			expectedCommitID: {MetaRangeID: expectedRangeID},
		},
	}
	// tests
	errSomethingBad := errors.New("first error")
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name: "without hook",
			hook: false,
			err:  nil,
		},
		{
			name: "hook no error",
			hook: true,
			err:  nil,
		},
		{
			name: "hook error",
			hook: true,
			err:  errSomethingBad,
		},
		{
			name:         "read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			expected := expectedCommitID
			refManager.TagCommitID = &expected
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreDeleteTagHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			err := g.DeleteTag(ctx, repo, expectedTagID, graveler.WithForce(tt.readOnlyRepo))

			// verify we got an error
			if !errors.Is(err, tt.err) {
				t.Fatalf("Delete tag err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Delete Tag err=%v, expected HookAbortError", err)
			}

			// verify that calls made until the first error
			called := slices.Contains(h.Called, "PreDeleteTagHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Pre delete Tag hook h.Called=%v, expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repository.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repository.RepositoryID)
			}
			if h.TagID != expectedTagID {
				t.Errorf("Hook tag ID '%s', expected '%s'", h.TagID, expectedTagID)
			}
		})
	}
}

func TestGraveler_PreCreateBranchHook(t *testing.T) {
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const sourceCommitID = graveler.CommitID("sourceCommitID")
	const sourceBranchID = graveler.CommitID("sourceBranchID")
	const newBranchPrefix = "newBranch-"
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		Refs: map[graveler.Ref]*graveler.ResolvedRef{graveler.Ref(sourceBranchID): {
			Type: graveler.ReferenceTypeBranch,
			BranchRecord: graveler.BranchRecord{
				BranchID: graveler.BranchID(sourceBranchID),
				Branch: &graveler.Branch{
					CommitID:     sourceCommitID,
					StagingToken: "",
				},
			},
		}},
	}
	// tests
	errSomethingBad := errors.New("first error")
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name: "without hook",
			hook: false,
			err:  nil,
		},
		{
			name: "hook no error",
			hook: true,
			err:  nil,
		},
		{
			name: "hook error",
			hook: true,
			err:  errSomethingBad,
		},
		{
			name:         "read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreCreateBranchHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}

			// WA for CreateBranch fake logic
			refManager.Branch = nil
			refManager.Err = graveler.ErrBranchNotFound

			newBranch := newBranchPrefix + strconv.Itoa(i)
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			_, err := g.CreateBranch(ctx, repo, graveler.BranchID(newBranch), graveler.Ref(sourceBranchID), graveler.WithForce(tt.readOnlyRepo))

			// verify we got an error
			if !errors.Is(err, tt.err) {
				t.Fatalf("Create branch err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Create branch err=%v, expected HookAbortError", err)
			}

			// verify that calls made until the first error
			called := slices.Contains(h.Called, "PreCreateBranchHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Pre-create branch hook h.Called=%v, expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repository.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repository.RepositoryID)
			}
			if h.CommitID != sourceCommitID {
				t.Errorf("Hook commit ID '%s', expected '%s'", h.BranchID, sourceCommitID)
			}
			if h.BranchID != graveler.BranchID(newBranch) {
				t.Errorf("Hook branch ID '%s', expected '%s'", h.BranchID, newBranch)
			}
		})
	}
}

func TestGraveler_PreDeleteBranchHook(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const sourceCommitID = graveler.CommitID("sourceCommitID")
	const sourceBranchID = graveler.CommitID("sourceBranchID")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	values := testutil.NewValueIteratorFake([]graveler.ValueRecord{{Key: nil, Value: nil}})
	stagingManager := &testutil.StagingFake{ValueIterator: values}
	refManager := &testutil.RefsFake{
		Refs: map[graveler.Ref]*graveler.ResolvedRef{graveler.Ref(sourceBranchID): {
			Type:                   graveler.ReferenceTypeBranch,
			ResolvedBranchModifier: 0,
			BranchRecord: graveler.BranchRecord{
				BranchID: graveler.BranchID(sourceBranchID),
				Branch: &graveler.Branch{
					CommitID:     sourceCommitID,
					StagingToken: "token",
				},
			},
		}},
		Branch:       &graveler.Branch{CommitID: sourceBranchID, StagingToken: "token"},
		StagingToken: "token",
	}
	// tests
	errSomethingBad := errors.New("first error")
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{
			name: "without hook",
			hook: false,
			err:  nil,
		},
		{
			name: "hook no error",
			hook: true,
			err:  nil,
		},
		{
			name: "hook error",
			hook: true,
			err:  errSomethingBad,
		},
		{
			name:         "read only repo",
			hook:         true,
			err:          nil,
			readOnlyRepo: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreDeleteBranchHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			err := g.DeleteBranch(ctx, repo, graveler.BranchID(sourceBranchID), graveler.WithForce(tt.readOnlyRepo))

			// verify we got an error
			if !errors.Is(err, tt.err) {
				t.Fatalf("Delete branch err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("Delete branch err=%v, expected HookAbortError", err)
			}

			// verify that calls made until the first error
			called := slices.Contains(h.Called, "PreDeleteBranchHook")
			if (tt.hook && !tt.readOnlyRepo) != called {
				t.Fatalf("Pre-delete branch hook h.Called=%v, expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}
			if !called {
				return
			}
			if h.RepositoryID != repository.RepositoryID {
				t.Errorf("Hook repository '%s', expected '%s'", h.RepositoryID, repository.RepositoryID)
			}
			if h.BranchID != graveler.BranchID(sourceBranchID) {
				t.Errorf("Hook branch ID '%s', expected '%s'", h.BranchID, sourceBranchID)
			}
		})
	}
}

func TestGravelerCreateCommitRecord(t *testing.T) {
	ctx := context.Background()
	t.Run("create commit record", func(t *testing.T) {
		test := testutil.InitGravelerTest(t)
		commit := graveler.Commit{
			Committer:    "committer",
			Message:      "message",
			MetaRangeID:  "metaRangeID",
			Parents:      []graveler.CommitID{"parent1", "parent2"},
			Metadata:     graveler.Metadata{"key": "value"},
			CreationDate: time.Now(),
			Version:      graveler.CurrentCommitVersion,
			Generation:   1,
		}
		test.RefManager.EXPECT().CreateCommitRecord(ctx, repository, graveler.CommitID("commitID"), commit).Return(nil)
		err := test.Sut.CreateCommitRecord(ctx, repository, "commitID", commit)
		require.NoError(t, err)
	})
}

func TestGraveler_Revert(t *testing.T) {
	type deps struct {
		CommittedManager *testutil.CommittedFake
		RefManager       *testutil.RefsFake
		StagingManager   *testutil.StagingFake
	}
	type args struct {
		branchID     graveler.BranchID
		ref          graveler.Ref
		parentNumber int
		allowEmpty   bool
	}
	tests := []struct {
		name        string
		deps        deps
		revertArgs  args
		expectedVal graveler.CommitID
		expectedErr error
	}{
		{
			name: "ref not found",
			deps: deps{
				CommittedManager: &testutil.CommittedFake{},
				RefManager: &testutil.RefsFake{
					Branch: &graveler.Branch{CommitID: "c1"},
					Commits: map[graveler.CommitID]*graveler.Commit{
						"c1": {MetaRangeID: "mri1"},
					},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
						"c1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				},
				StagingManager: &testutil.StagingFake{},
			},
			revertArgs: args{
				branchID:     "b1",
				ref:          graveler.Ref("ref1"),
				parentNumber: 0,
				allowEmpty:   false,
			},
			expectedErr: graveler.ErrCommitNotFound,
			expectedVal: "",
		},
		{
			name: "fail on staging token",
			deps: deps{
				CommittedManager: &testutil.CommittedFake{},
				RefManager: &testutil.RefsFake{
					Branch: &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{
						"c1": {MetaRangeID: "mri1"},
					},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"dirty-b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								BranchID: "dirty-b1",
								Branch: &graveler.Branch{
									CommitID:     "c1",
									StagingToken: "token",
								},
							},
						},
						"c1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				},
				StagingManager: &testutil.StagingFake{Value: value1, ValueIterator: testutil.NewValueIteratorFake(
					[]graveler.ValueRecord{{
						key1,
						value1,
					}})},
			},
			revertArgs: args{
				branchID:     "b1",
				ref:          graveler.Ref("c1"),
				parentNumber: 0,
				allowEmpty:   false,
			},
			expectedErr: graveler.ErrDirtyBranch,
			expectedVal: "",
		},
		{
			name: "fail on compacted data",
			deps: deps{
				CommittedManager: &testutil.CommittedFake{DiffIterator: testutil.NewDiffIter([]graveler.Diff{{Key: key1, Type: graveler.DiffTypeRemoved}})},
				RefManager: &testutil.RefsFake{
					Branch: &graveler.Branch{CommitID: "c1", StagingToken: "token", CompactedBaseMetaRangeID: "mri2"},
					Commits: map[graveler.CommitID]*graveler.Commit{
						"c1": {MetaRangeID: "mri1"},
					},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"dirty-b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								BranchID: "dirty-b1",
								Branch: &graveler.Branch{
									CommitID:                 "c1",
									StagingToken:             "token",
									CompactedBaseMetaRangeID: "mri2",
								},
							},
						},
						"c1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				},
				StagingManager: &testutil.StagingFake{Value: value1, ValueIterator: testutil.NewValueIteratorFake(
					[]graveler.ValueRecord{})},
			},
			revertArgs: args{
				branchID:     "b1",
				ref:          graveler.Ref("c1"),
				parentNumber: 0,
				allowEmpty:   false,
			},
			expectedErr: graveler.ErrDirtyBranch,
			expectedVal: "",
		},
		{
			name: "valid revert",
			deps: deps{
				CommittedManager: &testutil.CommittedFake{},
				RefManager: &testutil.RefsFake{
					Branch: &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{
						"c1": {MetaRangeID: "mri1"},
						"c2": {MetaRangeID: "mri2", Parents: graveler.CommitParents{"c1"}},
					},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID:     "c1",
									StagingToken: "token",
								},
							},
						},
						"c2": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c2",
								},
							},
						},
						"c1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				},
				StagingManager: &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(
					[]graveler.ValueRecord{{}})},
			},
			revertArgs: args{
				branchID:     "b1",
				ref:          graveler.Ref("c2"),
				parentNumber: 0,
				allowEmpty:   false,
			},
			expectedErr: nil,
			expectedVal: "",
		},
		{
			name: "valid revert and allow empty",
			deps: deps{
				CommittedManager: &testutil.CommittedFake{},
				RefManager: &testutil.RefsFake{
					Branch: &graveler.Branch{CommitID: "c1", StagingToken: "token"},
					Commits: map[graveler.CommitID]*graveler.Commit{
						"c1": {MetaRangeID: "mri1"},
						"c2": {MetaRangeID: "mri2", Parents: graveler.CommitParents{"c1"}},
					},
					Refs: map[graveler.Ref]*graveler.ResolvedRef{
						"b1": {
							Type:                   graveler.ReferenceTypeBranch,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								BranchID: "b1",
								Branch: &graveler.Branch{
									CommitID:     "c1",
									StagingToken: "token",
								},
							},
						},
						"c2": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c2",
								},
							},
						},
						"c1": {
							Type:                   graveler.ReferenceTypeCommit,
							ResolvedBranchModifier: 0,
							BranchRecord: graveler.BranchRecord{
								Branch: &graveler.Branch{
									CommitID: "c1",
								},
							},
						},
					},
				},
				StagingManager: &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(
					[]graveler.ValueRecord{{}})},
			},
			revertArgs: args{
				branchID:     "b1",
				ref:          graveler.Ref("c2"),
				parentNumber: 0,
				allowEmpty:   true,
			},
			expectedErr: nil,
			expectedVal: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGraveler(t, tt.deps.CommittedManager, tt.deps.StagingManager, tt.deps.RefManager, nil, testutil.NewProtectedBranchesManagerFake())

			got, err := g.Revert(context.Background(), repository, tt.revertArgs.branchID, tt.revertArgs.ref, tt.revertArgs.parentNumber, graveler.CommitParams{
				AllowEmpty: tt.revertArgs.allowEmpty,
			}, nil)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("unexpected err got = %v, wanted = %v", err, tt.expectedErr)
			}
			if err != nil {
				return
			}

			if got != tt.expectedVal {
				t.Errorf("got wrong commitID, got = %v, want %v", got, tt.expectedVal)
			}
		})
	}
}

func TestGraveler_CherryPickHooks(t *testing.T) {
	// prepare graveler
	const expectedRangeID = graveler.MetaRangeID("expectedRangeID")
	const expectedCommitID = graveler.CommitID("expectedCommitID")
	const sourceCommitID = graveler.CommitID("sourceCommitID")
	const cherryPickBranchID = graveler.BranchID("cherryPickBranchID")
	committedManager := &testutil.CommittedFake{MetaRangeID: expectedRangeID}
	stagingManager := &testutil.StagingFake{ValueIterator: testutil.NewValueIteratorFake(nil)}
	refManager := &testutil.RefsFake{
		CommitID: expectedCommitID,
		Branch:   &graveler.Branch{CommitID: expectedCommitID, StagingToken: "st1"},
		Commits: map[graveler.CommitID]*graveler.Commit{
			expectedCommitID: {MetaRangeID: expectedRangeID},
			sourceCommitID:   {MetaRangeID: expectedRangeID, Parents: graveler.CommitParents{expectedCommitID}},
		},
		Refs: map[graveler.Ref]*graveler.ResolvedRef{
			graveler.Ref(sourceCommitID): {
				Type: graveler.ReferenceTypeCommit,
				BranchRecord: graveler.BranchRecord{
					Branch: &graveler.Branch{
						CommitID: sourceCommitID,
					},
				},
			},
		},
	}
	// tests
	errSomethingBad := errors.New("cherry pick hook error")
	const commitCommitter = "committer"
	const cherryPickMessage = "cherry pick message"
	cherryPickMetadata := graveler.Metadata{"key1": "val1"}
	tests := []struct {
		name         string
		hook         bool
		err          error
		readOnlyRepo bool
	}{
		{name: "without hook", hook: false, readOnlyRepo: false},
		{name: "hook no error", hook: true, readOnlyRepo: false},
		{name: "hook read only repo", hook: true, readOnlyRepo: true},
		{name: "hook error", hook: true, err: errSomethingBad, readOnlyRepo: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			ctx := context.Background()
			g := newGraveler(t, committedManager, stagingManager, refManager, nil, testutil.NewProtectedBranchesManagerFake())
			h := &Hooks{Errs: map[string]error{
				"PreCherryPickHook":  tt.err,
				"PostCherryPickHook": tt.err,
			}}
			if tt.hook {
				g.SetHooksHandler(h)
			}
			repo := repository
			if tt.readOnlyRepo {
				repo = repositoryRO
			}
			// call cherry pick
			cherryPickCommitID, err := g.CherryPick(ctx, repo, cherryPickBranchID, graveler.Ref(sourceCommitID), nil, commitCommitter, &graveler.CommitOverrides{
				Message:  cherryPickMessage,
				Metadata: cherryPickMetadata,
			}, graveler.WithForce(tt.readOnlyRepo))

			// check err composition
			if !errors.Is(err, tt.err) {
				t.Fatalf("CherryPick err=%v, expected=%v", err, tt.err)
			}
			var hookErr *graveler.HookAbortError
			if err != nil && !errors.As(err, &hookErr) {
				t.Fatalf("CherryPick err=%v, expected HookAbortError", err)
			}

			// verify that calls made until the first error
			preCommitCalled := slices.Contains(h.Called, "PreCherryPickHook")
			postCommitCalled := slices.Contains(h.Called, "PostCherryPickHook")

			if (tt.hook && !tt.readOnlyRepo) != preCommitCalled {
				t.Fatalf("CherryPick invalid pre-cherry-pick hook call, %v expected=%t", h.Called, tt.hook && !tt.readOnlyRepo)
			}

			// PostCherryPick should only be called if PreCherryPick succeeded and operation completed
			if tt.hook && !tt.readOnlyRepo && tt.err == nil {
				if !postCommitCalled {
					t.Fatalf("CherryPick post-cherry-pick hook should be called when pre-cherry-pick succeeds, %v", h.Called)
				}
			} else {
				if postCommitCalled {
					t.Fatalf("CherryPick post-cherry-pick hook should not be called when pre-cherry-pick fails or repo is read-only, %v", h.Called)
				}
			}

			if !preCommitCalled {
				return
			}

			// verify hook parameters
			require.Equal(t, repo.RepositoryID, h.RepositoryID, "Hook repository ID mismatch")
			require.Equal(t, cherryPickBranchID, h.BranchID, "Hook branch ID mismatch")
			require.Equal(t, cherryPickMessage, h.Commit.Message, "Hook commit message mismatch")
			require.Equal(t, commitCommitter, h.Commit.Committer, "Hook commit committer mismatch")
			diff := deep.Equal(h.Commit.Metadata, cherryPickMetadata)
			require.Nil(t, diff, "Hook commit metadata diff:", diff)

			// verify post-cherry-pick hook parameters if it was called
			if postCommitCalled {
				require.Equal(t, cherryPickCommitID, h.CommitID, "Post-cherry-pick hook commit ID mismatch")
			}
		})
	}
}
