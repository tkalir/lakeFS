import React, {useContext, useEffect, useState} from "react";
import {useOutletContext} from "react-router-dom";

import {ActionGroup, ActionsBar, AlertError, Loading} from "../../../../lib/components/controls";
import {useRefs} from "../../../../lib/hooks/repo";
import {RepoError} from "../error";
import {useRouter} from "../../../../lib/hooks/router";
import Button from "react-bootstrap/Button";
import CompareBranches from "../../../../lib/components/repository/compareBranches";
import {RefTypeBranch} from "../../../../constants";
import Form from "react-bootstrap/Form";
import {pulls as pullsAPI} from "../../../../lib/api";
import CompareBranchesSelection from "../../../../lib/components/repository/compareBranchesSelection";
import {DiffContext, WithDiffContext} from "../../../../lib/hooks/diffContext";

const CreatePullForm = ({repo, reference, compare, title, setTitle, description, setDescription}) => {
    const router = useRouter();
    let [loading, setLoading] = useState(false);
    let [error, setError] = useState(null);

    const {state: {results: diffResults, loading: diffLoading, error: diffError}} = useContext(DiffContext);
    const isEmptyDiff = (!diffLoading && !diffError && !!diffResults && diffResults.length === 0);

    const onTitleInput = ({target: {value}}) => setTitle(value);
    const onDescriptionInput = ({target: {value}}) => setDescription(value);

    const submitForm = async () => {
        setLoading(true);
        setError(null);
        try {
            const {id: createdPullId} = await pullsAPI.create(repo.id, {
                title,
                description,
                source_branch: compare.id,
                destination_branch: reference.id
            });

            router.push({
                pathname: `/repositories/:repoId/pulls/:pullId`,
                params: {repoId: repo.id, pullId: createdPullId},
            });
        } catch (error) {
            setError(error.message);
            setLoading(false);
        }
    }

    return <>
        <Form.Group className="mb-3">
            <Form.Control
                required
                disabled={loading}
                type="text"
                size="lg"
                placeholder="Add a title..."
                value={title}
                onChange={onTitleInput}
            />
        </Form.Group>
        <Form.Group className="mb-3">
            <Form.Control
                required
                disabled={loading}
                as="textarea"
                rows={8}
                placeholder="Describe your changes..."
                value={description}
                onChange={onDescriptionInput}
            />
        </Form.Group>
        {error && <AlertError error={error} onDismiss={() => setError(null)}/>}
        <div>
            <Button variant="success"
                    disabled={!title || loading || isEmptyDiff}
                    onClick={submitForm}>
                {loading && <><span className="spinner-border spinner-border-sm text-light" role="status"/> {""}</>}
                Create Pull Request
            </Button>
            {isEmptyDiff &&
                <span className="alert alert-warning align-middle ms-4 pt-2 pb-2">
                    Pull requests must include changes.
                </span>
            }
        </div>
    </>;
};

const CreatePull = () => {
    const {repo, loading, error, reference, compare} = useRefs();

    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");

    if (loading) return <Loading/>;
    if (error) return <RepoError error={error}/>;

    return (
        <div className="w-75">
            <ActionsBar>
                <ActionGroup orientation="left">
                    <CompareBranchesSelection
                        repo={repo}
                        reference={reference}
                        compareReference={compare}
                        baseSelectURL={"/repositories/:repoId/pulls/create"}
                    />
                </ActionGroup>
            </ActionsBar>
            <h1 className="mt-3">Create Pull Request</h1>
            <div className="mt-4">
                <CreatePullForm
                    repo={repo}
                    reference={reference}
                    compare={compare}
                    title={title}
                    setTitle={setTitle}
                    description={description}
                    setDescription={setDescription}
                />
            </div>
            <hr className="mt-5 mb-4"/>
            <CompareBranches
                repo={repo}
                reference={{id: reference.id, type: RefTypeBranch}}
                compareReference={{id: compare.id, type: RefTypeBranch}}
            />
        </div>
    );
};

const RepositoryCreatePullPage = () => {
    const [setActivePage] = useOutletContext();
    useEffect(() => setActivePage("pulls"), [setActivePage]);
    return <WithDiffContext>
        <CreatePull/>
    </WithDiffContext>;
}

export default RepositoryCreatePullPage;
