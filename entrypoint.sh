#!/bin/bash

function parseInputs {
    subcommand="check"
    if [[ "${INPUT_FORCE_UPDATE}" == "true" ]]; then
        subcommand="update latest"
    elif [[ -n "${INPUT_SPECIFIC_VERSION}" ]]; then
        subcommand="update ${INPUT_SPECIFIC_VERSION}"
    fi

    workdir="./"
    if [[ -n "${INPUT_WORKING_DIR}" ]]; then
        workdir=${INPUT_WORKING_DIR}
    fi

    commentPR=""
    if [[ "${INPUT_COMMENT_PR}" == "true" ]]; then
        commentPR="true"
    fi
}

function main {
    parseInputs
    output=$(go run main.go ${subcommand} --root-path ${workdir} 2> /dev/null)
    exitCode=${?}

    if [ ${exitCode} -ne 0  &&  ${exitCode} -ne 3 ]; then
        echo "error: failed to check Terraform cloud version"
        exit ${exitCode}
    fi

    if [ -z "${output}" ]; then
        echo "::set-output name=is_available_update::false"
        echo "info: no updates are available"
        exit 0
    fi

    echo "::set-output name=output::${output}"
    echo "::set-output name=is_available_update::true"

    workspaceLink=$(echo "${output}" | tail -n 1)
    output=$(echo "${output}" | grep -v "${workspaceLink}")
    echo "info: ${output}"

    if [[ "${commentPR}" == "true" ]] && [[ "${GITHUB_EVENT_NAME}" == "pull_request" ]]; then
        CommentsURL=$(cat ${GITHUB_EVENT_PATH} | jq -r .pull_request.comments_url)
        echo "info: commenting on the pull request"
        output="{\"body\": \"#### Terraform Cloud Workspace new version available\n\`\`\`\n${output}\n\`\`\`\n\n*working directory: \`${workdir}\`, ${workspaceLink}*\"}"
        echo "${output}"  | curl -XPOST -sS -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: application/json" --data @- "${CommentsURL}"
        if [ ${?} -ne 0 ]; then
            echo "error: failed to post a comment to the pull request"
        fi
    fi
}

main
