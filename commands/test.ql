
query getFolder2($owner: String!, $name: String!,$branch: GitObjectID) {
  repository(owner: $owner, name: $name) {
    object(expression: "master:",oid: $branch) {
      ... on Tree {
        repository {
          id
          name
        }

        entries {
					oid
          name
          type
          repository {
            name
          }
        }
      }
    }
  }
}
query getFolder($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    object(expression: "master:") {
      ... on Tree {
        repository {
          id
          name
        }

        entries {
					oid
          name
          type
          repository {
            name
          }
        }
      }
    }
  }
}

query getListFolder($owner: String!, $per_page: Int = 100, $endCursor: String) {
  repositoryOwner(login: $owner) {
    repositories(first: $per_page, after: $endCursor, ownerAffiliations: OWNER) {
      nodes {
        nameWithOwner
        id
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
}

query getRepository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    id
    databaseId
  }
}

query getDraftPR($owner: String!, $name: String!, $headRef: String!) {
  repository(owner: $owner, name: $name) {
    pullRequests(headRefName: $headRef, first: 100) {
      nodes {
        id
        state
      }
    }
  }
}

query getIssue($owner: String!, $name: String!, $number: Int!) {
  repository(owner: $owner, name: $name) {
    issue(number: $number) {
      title
    }
  }
}

query getAllIssuesPullRequestSinceDate($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    object(expression: "master:") {
      ... on Tree {
        entries {
          name
          type
          mode
        }
      }
    }
  }
}

query getListFilesToCommit($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    ref(qualifiedName: "master") {
      target {
        ... on Commit {
          id
          history(first: 5) {
            pageInfo {
              hasNextPage
            }
            edges {
              node {
                changedFiles
                messageHeadline
                oid
                message
                author {
                  name
                  email
                  date
                }
                commitResourcePath
              }
            }
          }
        }
      }
    }
  }
}

//////////////

query filteredTree($owner: String!, $name: String!,$branch: GitObjectID) {
  repository(owner: $owner, name: $name) {
    object(expression: "master:",oid: $branch) {
      ... on Tree {
        entries {
					oid
          name
          type
        }
      }
    }
  }
}
query dirTree($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    object(expression: "master:") {
      ... on Tree {
        repository {
          id
          name
        }

        entries {
					oid
          name
          type
        }
      }
    }
  }
}

api --paginate graphql -f owner="hilalisadev" name="hilalisadev.github.io" branch="1c0e3a7cd3a64b26ab4b6afce7aab55bbcde48a0" "$@" -f query='
     query filteredTree($owner: String!, $name: String!,$branch: GitObjectID) {
     repository(owner: $owner, name: $name) {
       object(expression: "master:",oid: $branch) {
         ... on Tree {
           entries {
   					oid
             name
             type
           }
         }
       }
     }
   }
     '


go run .\main.go api --paginate graphql -f master="master:" owner="hilalisadev" name="hilalisadev.github.io" branch="1c0e3a7cd3a64b26ab4b6afce7aab55bbcde48a0" "$@" -f query='  query filteredTree($owner: String!, $name: String!,$branch: GitObjectID, $master: String!) {
  repository(owner: $owner, name: $name) {
    object(expression: $master,oid: $branch) {
      ... on Tree {
        entries {
                    oid
          name
          type
        }
      }
    }
  }
}
 '
>>
query filteredTree($owner: String!, $name: String!,$branch: GitObjectID) {
  repository(owner: $owner, name: $name) {
    object(expression: "master:",oid: $branch) {
      ... on Tree {
        entries {
					oid
          name
          type
        }
      }
    }
  }
}

>>
api
--paginate
graphql
-f
owner="hilalisadev"
-f
name="hilalisadev.github.io"
-f
branch="1c0e3a7cd3a64b26ab4b6afce7aab55bbcde48a0"
-f
master="master:"
"$@"
-f
query=" query filteredTree($owner: String!, $name: String!,$branch: GitObjectID, $master: String!) { repository(owner: $owner, name: $name) { object(expression: $master,oid: $branch) { ... on Tree { entries { oid name type } } } } } "


>>
{"data":{"repository":{"object":{"entries":[{"oid":"a592e6d3366a3ac0ae1b62f40ae205ceac984fa3","name":"1.js","type":"blob"},{"oid":"cadcb0aabe792db6a0f417567e3dbe842178286e","name":"2.js","type":"blob"},{"oid":"1bc45fb1fb832208c6b3d775e19ceb245fa2a0a9","name":"3.js","type":"blob"},{"oid":"6cfc1462e4bd6fae2798d00d2f69cd5fd1db9089","name":"4.js","type":"blob"},{"oid":"4f8dd21150bd3556656bd2b3004304feb019e6e3","name":"images","type":"tree"},{"oid":"366e9c20672929459ba9d32d2375b69977157c81","name":"main.js","type":"blob"},{"oid":"c59ef7d371c2420b0ef6030016c7112dc137fd87","name":"manifest.json","type":"blob"},{"oid":"c9e1e585795fbe2928c7aad0b7c8da2740d3c82f","name":"node_modules","type":"tree"},{"oid":"b025ae1c982a0f82ab6ef845f1327ca12300a5f3","name":"precache-manifest.a122b2567c9493eb570019a9a840a929.js","type":"blob"},{"oid":"0ef14ea65883627156074ce6e44e96a8133fed60","name":"service-worker.js","type":"blob"}]}}}}
