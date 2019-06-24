#!/usr/bin/env bash

# Safe check - skips relese commit generation when already tagged commit
if [[ $(git name-rev --name-only --tags HEAD) = "v$CURRENT_VERSION" ]]; then
    echo "Already tagged or no new commits introduced. Skipping.."
    exit 0
fi

if [ "" == "$GH_USER_TOKEN" ]; then
    echo "GitHub Token is not provided."
    exit 1
fi

if [ "" == "$GH_USER_NAME" ]; then
    echo "GitHub User name is not provided."
    exit 1
fi

if [ "" == "$GH_USER_EMAIL" ]; then
    echo "GitHub User email is not provided."
    exit 1
fi

git config --global user.email "$GH_USER_EMAIL"
git config --global user.name "$GH_USER_NAME"

# Guess new version number
RECOMMENDED_BUMP=$(conventional-recommended-bump -p angular)

# Split version by dots
IFS='.' read -r -a V <<< "$CURRENT_VERSION"

# Ignore postfix like "-dev"
((V[2]++))
((V[2]--))
CURRENT_VERSION_SEM="${V[0]}.${V[1]}.${V[2]}"

# When version is 0.x.x it is allowed to make braking changes on minor version
if [[ "0" = "${V[0]}" ]] && [[ "${RECOMMENDED_BUMP}" = "major" ]]; then
    RECOMMENDED_BUMP="minor";
fi

echo "Recommended bump: $RECOMMENDED_BUMP"

# Increment semantic version numbers major.minor.patch
if [[ "${RECOMMENDED_BUMP}" = "major" ]]; then
    ((V[0]++));
    V[1]=0;
    V[2]=0;
elif [[ "${RECOMMENDED_BUMP}" = "minor" ]]; then
    ((V[1]++));
    V[2]=0;
elif [[ "${RECOMMENDED_BUMP}" = "patch" ]]; then ((V[2]++));
else
    echo "Could not bump version"
    exit 1
fi

NEW_VERSION_SEM="${V[0]}.${V[1]}.${V[2]}"
NEW_VERSION=${CURRENT_VERSION//${CURRENT_VERSION_SEM}/${NEW_VERSION_SEM}}

echo "Old version: ${CURRENT_VERSION} ($CURRENT_VERSION_SEM)"
echo "New version: ${NEW_VERSION} ($NEW_VERSION_SEM)"

# Tag to update changelog with new version included
git tag "v${NEW_VERSION}"
conventional-changelog -p angular -i CHANGELOG.md -s -r 2
git tag -d "v${NEW_VERSION}"

# Create release commit and tag
git add CHANGELOG.md
git commit -m "chore(release): v${NEW_VERSION} :tada:
$(conventional-changelog)
"
git tag "v${NEW_VERSION}"

if [ "" == "$GH_USER" ]; then
    GH_USER="$CIRCLE_PROJECT_USERNAME";
fi

# Push commit and tag
git remote add authorized "https://${GH_USER_NAME}:${GH_USER_TOKEN}@github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}.git"
git push authorized HEAD:master --tags
git push authorized HEAD:develop

# Make github release
conventional-github-releaser -p angular -t "${GH_USER_TOKEN}"
