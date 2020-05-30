#!/bin/bash

function parseInputs {
    subcommand="check"
    if [[ "${INPUT_FORCE_UPDATE}" == "true" ]]; then
        subcommand="update latest"
    fi

    tfetoken=""
    if [[ -n "${INPUT_TFE_TOKEN}" ]]; then
        tfetoken=${INPUT_TFE_TOKEN}
    fi

    workdir="./"
    if [[ -n "${INPUT_WORKDIR}" ]]; then
        workdir=${INPUT_WORKDIR}
    fi

    commentPR=""
    if [[ "${INPUT_COMMENT_PR}" == "true" ]]; then
        commentPR="true"
    fi
}

function main {
    parseInputs
    output=$(go run main.go ${subcommand} --token ${tfetoken} --root-path ${workdir}) 2> /dev/null
    exitCode=${?}

    if [ ${exitCode} -ne 0 ]; then
        echo "error: failed to check Terraform cloud version"
        exit ${exitCode}
    fi

    if [ -z "${output}" ]; then
        echo "info: no updates are available"
        exit 0
    fi

    echo "::set-output name=result::${output}"

    if [[ "${commentPR}" == "true" ]] && [[ "${GITHUB_EVENT_NAME}" == "pull_request" ]]; then
        CommentsURL=$(cat ${GITHUB_EVENT_PATH} | jq -r .pull_request.comments_url)
        echo "info: commenting on the pull request"
        output="Terraform Cloud Workspace [[${subcommand}]] has detected new version\n\n${output}"
        echo "${output}" | curl -s -S -H "Authorization: token ${GITHUB_TOKEN}" --header "Content-Type: application/json" --data @- "${CommentsURL}" > /dev/null
    fi
}

main
