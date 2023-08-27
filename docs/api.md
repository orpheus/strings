# API Docs

## Threads & Strings

### POST /v1/threads

#### Rules

Thread

1. Posting empty body results in error, missing `name`
2. `name` required field (`name` is not unique in db (different versions share same name))
3. `deleted`, `archived`, `version`, and `dateCreated` are client-immutable fields (only server can change)
4. `id` will always be generated by server

String

1. Posting empty body results in error, missing `name`
2. `name` required field (`name` is not unique in db (different versions share same name))
3. `order` required field
4. `deleted`, `archived`, `version`, `active`, and `dateCreated` are client-immutable fields (only server can change)
5. `id` will always be generated by server

#### Workflows and requirements

1. Create new thread
   1. required: `name`
      1. `id` will be ignored (always generated by the server)
      2. `threadId` exists
         1. if thread exists, will update existing thread (see 2)
         2. if thread not exists, `threadId` be used as the server's `threadId`
      3. `threadId` not exists
         1. server will generate one and create a new thread (a thread record and a versioned_thread record with version:1)
2. Update existing thread
   1. required: `threadId` (if field not exists, will create new thread -- see #1)
      1. `id` will be ignored
      2. if thread not exists, will create new thread (see 1)
      3. if thread exists, will update existing thread
         1. if thread changed, update thread
            1. create new thread version with bumped version
            2. overwrite any client-changeable fields
               1. name
               2. strings
                  1. TODO: for strings
                  2. if thread not changed, return server fetched thread (each post will always return the full thread, even if I just supply the `threadId` and `name`)
3. Create new string
   1. required: `name`
   2. optional: `order`
      1. exists
         1. validate `order` is valid
            1. not beyond known limit
      2. not exists
         1. `order` will be assigned to last in limit
   3. will always use `threadId` defined by the thread, client-supplied id will be ignored
   4. a new `stringId` will be created for this string associated with the `threadId` of the parent thread
   5. `id` is ignored
4. Update string
   1. required: `name`
   2. required: `stringId` must be valid
      1. if `stringId` not found, a new versioned string will be created for that thread
   3. `id` is ignored
5. Creating/Updating multiple strings
   1. required: `thread_id` uuid
   2. required: `strings` []string
      1. grab server thread with server strings
      2. update server strings from client strings (client-update fields only)
      3. determine/validate order between existing strings and new strings after existing strings have been updated
         1. server always assumes client provided Perfect Order. server will not ever try to order strings for you.
      4. set `order` on new strings without `order` to last in range
      5. create new versions for updated strings
      6. create new string records for new strings
      7. return and bump thread version

#### Cases

##### Detecting object change

1. if the same request payload is POSTd multiple times, want to prevent unnecessary object version creations
2. need to check if content is same/diff
   1. hash only relevant CONTENT, not metadata
      1. `name`
      2. `strings`
         1. same applies for strings, want to ignore metadata content

Solution

Idea 1

Run a diff compare on each object where it checks only the relevant content

1. Get thread and all strings from database
   1. get latest core.Thread object and compare to POSTd core.Thread object
2. Each thread/string object has a .compare() method which reads in only relevant content data
3. Run through each string running .compare(), updating each one that diffs
4. Diff the thread and update if needed
5. If any string diffs, bump thread

#### Examples

##### Create new thread

Request

```bash
curl --location --request POST '0.0.0.0:8080/v1/threads' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "ryan"
}'
```

Response

```json
{
   "id": "1cad1b42-34c9-4d77-a5c1-33f3fe0855d9",
   "name": "ryan",
   "version": 1,
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101",
   "archived": false,
   "deleted": false,
   "dateCreated": "2023-08-26T19:03:00.431708-07:00",
   "strings": null
}
```

##### POST existing thread

Request

Note here the `thread_id` is the same as the one above

```bash
curl --location --request POST '0.0.0.0:8080/v1/threads' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "ryan",
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101"
}'
```

Response

Payload is the same as the first POST where the thread was created

```json
{
   "id": "1cad1b42-34c9-4d77-a5c1-33f3fe0855d9",
   "name": "ryan",
   "version": 1,
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101",
   "archived": false,
   "deleted": false,
   "dateCreated": "2023-08-26T19:03:00.431708-07:00",
   "strings": null
}
```

##### Update existing thread

Request

1. Updated `name` from `ryan` to `ryan-2`

```bash
curl --location --request POST '0.0.0.0:8080/v1/threads' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "ryan-2",
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101"
}'
```

Response

1. name is updated
2. version bumped
3. new id was generated

```json
{
   "id": "9bfb6090-0a35-48e3-b39e-9eb5b3de020c",
   "name": "ryan-2",
   "version": 2,
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101",
   "archived": false,
   "deleted": false,
   "dateCreated": "2023-08-26T19:07:15.219495-07:00",
   "strings": null
}
```

##### Attempt to update server controlled thread fields

Request

1. `name` and `thread_id` are the same
2. set `deleted` and `archived` (not allowed for client to update)

```bash
curl --location --request POST '0.0.0.0:8080/v1/threads' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name": "ryan-2",
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101",
   "deleted": true,
   "archived": true
}'
```

Response

1. `archived` and `deleted` are not set by client

```json
{
   "id": "af5f6460-d4bd-4c9f-8edd-d1c1366e3c87",
   "name": "ryan-2",
   "version": 1,
   "thread_id": "ee8c1f41-d3dc-44f4-8200-1dc89a2d3101",
   "archived": false,
   "deleted": false,
   "dateCreated": "2023-08-26T19:52:13.492728-07:00",
   "strings": null
}
```