---
title: Garbage Collection
description: Clean up expired objects using the garbage collection feature in lakeFS.
---

# Garbage Collection

!!! tip
    [lakeFS Cloud](https://lakefs.cloud) users enjoy a [managed garbage collection](./managed-gc.md) service, and do not need to run this Spark program.

!!! tip
    [lakeFS Enterprise](../../enterprise/index.md) users can run a [stand alone GC program](./standalone-gc.md), instead of this Spark program.

By default, lakeFS keeps all your objects forever. This allows you to travel back in time to previous versions of your data.
However, sometimes you may want to remove the objects from the underlying storage completely.
Reasons for this include cost-reduction and privacy policies.

The garbage collection (GC) job is a Spark program that removes the following from the underlying storage:

1. _Committed objects_ that have been deleted (or replaced) in lakeFS, and are considered expired according to [rules you define](#garbage-collection-rules).
2. _Uncommitted objects_ that are no longer accessible
    * For example, objects deleted before ever being committed.

## Garbage collection rules

!!! info
    These rules only apply to objects that have been _committed_ at some point.
    Without retention rules, only inaccessible _uncommitted_ objects will be removed by the job.

Garbage collection rules determine for how long an object is kept in the storage after it is _deleted_ (or replaced) in lakeFS.
For every branch, the GC job retains deleted objects for the number of days defined for the branch.
In the absence of a branch-specific rule, the default rule for the repository is used.
If an object is present in more than one branch ancestry, it is removed only after the retention period has ended for
all relevant branches.

Example GC rules for a repository:

```json
{
  "default_retention_days": 14,
  "branches": [
    {"branch_id": "main", "retention_days": 21},
    {"branch_id": "dev", "retention_days": 7}
  ]
}
```

In the above example, objects will be retained for 14 days after deletion by default.
However, if present in the branch `main`, objects will be retained for 21 days.
Objects present _only_ in the `dev` branch will be retained for 7 days after they are deleted.

### How to configure garbage collection rules

To define retention rules, either use the `lakectl` command, the lakeFS web UI, or [API](../../reference/api.md#/retention/set%20garbage%20collection%20rules):

=== "CLI"
    Create a JSON file with your GC rules:

    ```bash
    cat <<EOT >> example_repo_gc_rules.json
    {
    "default_retention_days": 14,
    "branches": [
        {"branch_id": "main", "retention_days": 21},
        {"branch_id": "dev", "retention_days": 7}
    ]
    }
    EOT
    ```

    Set the GC rules using `lakectl`:
    ```bash
    lakectl gc set-config lakefs://example-repo -f example_repo_gc_rules.json 
    ```
=== "Web UI"
    From the lakeFS web UI:

    1. Navigate to the main page of your repository.
    2. Go to _Settings_ -> _Garbage Collection_.
    3. Click _Edit policy_ and paste your GC rule into the text box as a JSON.
    4. Save your changes.

    ![GC Rules From UI](../../assets/img/gc_rules_from_ui.png)

## How to run the garbage collection job

To run the job, use the following `spark-submit` command (or using your preferred method of running Spark programs).

=== "AWS"
    ```bash
    spark-submit --class io.treeverse.gc.GarbageCollection \
        --packages org.apache.hadoop:hadoop-aws:2.7.7 \
        -c spark.hadoop.lakefs.api.url=https://lakefs.example.com:8000/api/v1  \
        -c spark.hadoop.lakefs.api.access_key=<LAKEFS_ACCESS_KEY> \
        -c spark.hadoop.lakefs.api.secret_key=<LAKEFS_SECRET_KEY> \
        -c spark.hadoop.fs.s3a.access.key=<S3_ACCESS_KEY> \
        -c spark.hadoop.fs.s3a.secret.key=<S3_SECRET_KEY> \
        http://treeverse-clients-us-east.s3-website-us-east-1.amazonaws.com/lakefs-spark-client/0.15.0/lakefs-spark-client-assembly-0.15.0.jar \
        example-repo us-east-1
    ```

=== "Azure"
    If you want to access your storage using the account key:

    ```bash
    spark-submit --class io.treeverse.gc.GarbageCollection \
        --packages org.apache.hadoop:hadoop-aws:3.2.1 \
        -c spark.hadoop.lakefs.api.url=https://lakefs.example.com:8000/api/v1  \
        -c spark.hadoop.lakefs.api.access_key=<LAKEFS_ACCESS_KEY> \
        -c spark.hadoop.lakefs.api.secret_key=<LAKEFS_SECRET_KEY> \
        -c spark.hadoop.fs.azure.account.key.<AZURE_STORAGE_ACCOUNT>.dfs.core.windows.net=<AZURE_STORAGE_ACCESS_KEY> \
        http://treeverse-clients-us-east.s3-website-us-east-1.amazonaws.com/lakefs-spark-client/0.15.0/lakefs-spark-client-assembly-0.15.0.jar \
        example-repo
    ```

    Or, if you want to access your storage using an Azure service principal:

    ```bash
    spark-submit --class io.treeverse.gc.GarbageCollection \
        --packages org.apache.hadoop:hadoop-aws:3.2.1 \
        -c spark.hadoop.lakefs.api.url=https://lakefs.example.com:8000/api/v1  \
        -c spark.hadoop.lakefs.api.access_key=<LAKEFS_ACCESS_KEY> \
        -c spark.hadoop.lakefs.api.secret_key=<LAKEFS_SECRET_KEY> \
        -c spark.hadoop.fs.azure.account.auth.type.<AZURE_STORAGE_ACCOUNT>.dfs.core.windows.net=OAuth \
        -c spark.hadoop.fs.azure.account.oauth.provider.type.<AZURE_STORAGE_ACCOUNT>.dfs.core.windows.net=org.apache.hadoop.fs.azurebfs.oauth2.ClientCredsTokenProvider \
        -c spark.hadoop.fs.azure.account.oauth2.client.id.<AZURE_STORAGE_ACCOUNT>.dfs.core.windows.net=<application-id> \
        -c spark.hadoop.fs.azure.account.oauth2.client.secret.<AZURE_STORAGE_ACCOUNT>.dfs.core.windows.net=<service-credential-key> \
        -c spark.hadoop.fs.azure.account.oauth2.client.endpoint.<AZURE_STORAGE_ACCOUNT>.dfs.core.windows.net=https://login.microsoftonline.com/<directory-id>/oauth2/token \
        http://treeverse-clients-us-east.s3-website-us-east-1.amazonaws.com/lakefs-spark-client/0.15.0/lakefs-spark-client-assembly-0.15.0.jar \
        example-repo
    ```

    !!! note
        * On Azure, GC was tested only on Spark 3.3.0, but may work with other Spark and Hadoop versions.
        * In case you don't have `hadoop-azure` package as part of your environment, you should add the package to your spark-submit with `--packages org.apache.hadoop:hadoop-azure:3.2.1`
        * For GC to work on Azure blob, [soft delete](https://docs.microsoft.com/en-us/azure/storage/blobs/soft-delete-blob-overview) should be disabled.
        

=== "GCP"
    For Garbage Collection to work with GCP, you must provide it with a service account key JSON file. 
    The use service account must have `Storage Object User` permissions for the repository namespace (bucket).

    ```bash
    spark-submit --class  io.treeverse.gc.GarbageCollection \
        -c spark.hadoop.lakefs.api.url=https://lakefs.example.com:8000/api/v1  \
        -c spark.hadoop.lakefs.api.access_key=<LAKEFS_ACCESS_KEY> \
        -c spark.hadoop.lakefs.api.secret_key=<LAKEFS_SECRET_KEY> \
        -c spark.hadoop.google.cloud.auth.service.account.enable=true \
        -c spark.hadoop.google.cloud.auth.service.account.json.keyfile=<PATH_TO_JSON_KEYFILE> \
        -c spark.hadoop.fs.gs.impl=com.google.cloud.hadoop.fs.gcs.GoogleHadoopFileSystem \
        -c spark.hadoop.fs.AbstractFileSystem.gs.impl=com.google.cloud.hadoop.fs.gcs.GoogleHadoopFS \
        http://treeverse-clients-us-east.s3-website-us-east-1.amazonaws.com/lakefs-spark-client/0.15.0/lakefs-spark-client-assembly-0.15.0.jar \
        example-repo
    ```

### Mark and Sweep stages

You can break the job into two stages:

* _Mark_: find objects to remove, without actually removing them.
* _Sweep_: remove the objects.

#### Mark-only mode

To make GC run the mark stage only, add the following to your spark-submit command:

```properties
spark.hadoop.lakefs.gc.do_sweep=false
```

In mark-only mode, GC will write the keys of the expired objects under: `<REPOSITORY_STORAGE_NAMESPACE>/_lakefs/retention/gc/unified/<MARK_ID>/`.
_MARK_ID_ is generated by the job. You can find it in the driver's output:

```
Report for mark_id=gmc6523jatlleurvdm30 path=s3a://example-bucket/_lakefs/retention/gc/unified/gmc6523jatlleurvdm30
```

#### Sweep-only mode

To make GC run the sweep stage only, add the following properties to your spark-submit command:

```properties
spark.hadoop.lakefs.gc.do_mark=false
spark.hadoop.lakefs.gc.mark_id=<MARK_ID> # Replace <MARK_ID> with the identifier you obtained from a previous mark-only run
```

## Garbage collection notes

1. In order for an object to be removed, it must not exist on the HEAD of any branch.
   You should remove stale branches to prevent them from retaining old objects.
   For example, consider a branch that has been merged to `main` and has become stale.
   An object which is later deleted from `main` will always be present in the stale branch, preventing it from being removed.
1. lakeFS will never delete objects outside your repository's storage namespace.
   In particular, objects that were imported using `lakectl import` or the UI import wizard will not be affected by GC jobs.
1. In cases where deleted objects are brought back to life while a GC job is running (for example, by reverting a commit),
   the objects may or may not be deleted.
1. Garbage collection does not remove any commits: you will still be able to use commits containing removed objects,
   but trying to read these objects from lakeFS will result in a `410 Gone` HTTP status.
