echo machine github.com        >  .netrc
echo login $GITHUB_USERNAME    >> .netrc
echo password $GITHUB_TOKEN    >> .netrc
echo $GITHUB_USERNAME

docker build \
    --no-cache \
    -t scraper \
    --secret id=netrc,src="./.netrc" \
    .