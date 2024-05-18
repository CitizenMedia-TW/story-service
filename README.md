# story-service

- RESTful APIs running on `localhost:50051`

---

### MyStory

<details>
<summary><code>GET</code> <code><b>/mystory</b></code> <code>(Get stories written by the user)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Responses

> | http code    | content-type       | response                                          |
> | ------------ | ------------------ | ------------------------------------------------- |
> | `200`        | `application/json` | `{"message": "Success", "storyIdList": string[]}` |
> | `401`, `500` | `text/plain`       | N/A                                               |

</details>

---

### Story

<details>
<summary><code>POST</code> <code><b>/story</b></code> <code>(Create a story)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Body (application/json)

> | key      | required | data type | description            |
> | -------- | -------- | --------- | ---------------------- |
> | authorId | true     | string    | ObjectId of the author |
> | content  | true     | string    | Content of the story   |
> | title    | true     | string    | Title of the story     |
> | subtitle | true     | string    | Subtitle of the story  |
> | tags     | true     | string[]  | Tags of the story      |

##### Responses

> | http code           | content-type       | response                                                         |
> | ------------------- | ------------------ | ---------------------------------------------------------------- |
> | `200`               | `application/json` | `{"message": "Success", "storyId": "ObjectId of the new story"}` |
> | `401`, `400`, `500` | `text/plain`       | N/A                                                              |

</details>

<details>
<summary><code>GET</code> <code><b>/story</b></code> <code>(Get a story by given id)</code></summary>

##### Query Parameters

> | key     | required | data type | description           |
> | ------- | -------- | --------- | --------------------- |
> | storyId | true     | string    | ObjectId of the story |

##### Responses

```typescript
type story = {
  author: string
  authorId: string
  content: string
  title: string
  subTitle: string
  createdAt: google.protobuf.Timestamp
  comments: Comment
  tags: string[]
}
```

> | http code           | content-type       | response                                 |
> | ------------------- | ------------------ | ---------------------------------------- |
> | `200`               | `application/json` | `{"message": "Success", "story": story}` |
> | `400`, `404`, `500` | `text/plain`       | N/A                                      |

</details>

<details>
<summary><code>DELETE</code> <code><b>/story</b></code> <code>(Delete a story by given id)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Body (application/json)

> | key       | required | data type | description             |
> | --------- | -------- | --------- | ----------------------- |
> | storyId   | true     | string    | ObjectId of the story   |
> | deleterId | true     | string    | ObjectId of the deleter |

##### Responses

> | http code    | content-type       | response                 |
> | ------------ | ------------------ | ------------------------ |
> | `200`        | `application/json` | `{"message": "Success"}` |
> | `400`, `500` | `text/plain`       | N/A                      |

</details>

---

### Recommend

<details>
<summary><code>GET</code> <code><b>/recommend</b></code> <code>(Get resommended stories, WIP)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Query Parameters

> | key    | required | data type | description                 |
> | ------ | -------- | --------- | --------------------------- |
> | userId | true     | string    | ObjectId of the user        |
> | count  | true     | int       | Number of story to retrieve |
> | skip   | true     | int       | --                          |

##### Responses

> | http code    | content-type       | response                                          |
> | ------------ | ------------------ | ------------------------------------------------- |
> | `200`        | `application/json` | `{"message": "Success", "storyIdList": string[]}` |
> | `400`, `500` | `text/plain`       | N/A                                               |

</details>

---

### Comment

<details>
<summary><code>POST</code> <code><b>/comment</b></code> <code>(Add a new comment to the story)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Body (application/json)

> | key              | required | data type | description                |
> | ---------------- | -------- | --------- | -------------------------- |
> | comment          | true     | string    | The content of the comment |
> | commenterId      | true     | string    | The id of the commenter    |
> | commentedStoryId | true     | string    | The story to comment on    |

##### Responses

> | http code           | content-type       | response                                                  |
> | ------------------- | ------------------ | --------------------------------------------------------- |
> | `200`               | `application/json` | `{"message": "Success", "commentId: "id of the comment"}` |
> | `400`, `401`, `500` | `text/plain`       | N/A                                                       |

</details>

<details>
<summary><code>DELETE</code> <code><b>/comment</b></code> <code>(Delete a comment of a story)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Body (application/json)

> | key       | required | data type | description           |
> | --------- | -------- | --------- | --------------------- |
> | commentId | true     | string    | The id of the comment |
> | deleterId | true     | string    | The id of the deleter |

##### Responses

> | http code           | content-type       | response                 |
> | ------------------- | ------------------ | ------------------------ |
> | `200`               | `application/json` | `{"message": "Success"}` |
> | `400`, `401`, `500` | `text/plain`       | N/A                      |

</details>

---

### SubComment

<details>
<summary><code>POST</code> <code><b>/subComment</b></code> <code>(Add a new subComment to the comment)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Body (application/json)

> | key              | required | data type | description                |
> | ---------------- | -------- | --------- | -------------------------- |
> | content          | true     | string    | The content of the comment |
> | commenterId      | true     | string    | The id of the commenter    |
> | repliedCommentId | true     | string    | The comment to comment on  |

##### Responses

> | http code    | content-type       | response                                                        |
> | ------------ | ------------------ | --------------------------------------------------------------- |
> | `200`        | `application/json` | `{"message": "Success", "subCommentId: "id of the subComment"}` |
> | `401`, `500` | `text/plain`       | N/A                                                             |

</details>

<details>
<summary><code>DELETE</code> <code><b>/comment</b></code> <code>(Delete a subComment of a comment)</code></summary>

##### Headers

> | key           | value          | description   |
> | ------------- | -------------- | ------------- |
> | Authorization | `Bearer token` | The jwt token |

##### Body (application/json)

> | key       | required | data type | description           |
> | --------- | -------- | --------- | --------------------- |
> | commentId | true     | string    | The id of the comment |
> | deleterId | true     | string    | The id of the deleter |

##### Responses

> | http code          | content-type       | response                 |
> | ------------------ | ------------------ | ------------------------ |
> | `200`              | `application/json` | `{"message": "Success"}` |
> | `400`, `401`, 500` | `text/plain`       | N/A                      |

</details>

---

### Search

<details>
<summary><code>GET</code> <code><b>/search</b></code> <code>(Search for the stories)</code></summary>

##### Query Parameters

> | key   | required | data type | description                 |
> | ----- | -------- | --------- | --------------------------- |
> | tag   | true     | string    | Tag for searching           |
> | count | true     | int       | Number of story to retrieve |
> | skip  | true     | int       | --                          |

##### Responses

> | http code | content-type       | response                                                   |
> | --------- | ------------------ | ---------------------------------------------------------- |
> | `200`     | `application/json` | `{"message": "Success", "storyIdList: string list of ids}` |
> | `500`     | `text/plain`       | N/A                                                        |

</details>
