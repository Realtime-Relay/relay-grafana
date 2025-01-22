echo "Creating data source package..."

PLUGIN_VERSION=$(cat package.json | jq .version -r)

echo "Plugin Version => $PLUGIN_VERSION"

echo "Creating plugin binaries..."
webpack -c ./.config/webpack/webpack.config.ts --env production 
mage 

echo "Updating README..."
rm dist/README.md
cp README.md dist/README.md

echo "Creating release..."
mkdir releases/arjunpreetham49-relay-datasource
mv dist/* releases/arjunpreetham49-relay-datasource

echo "Moving GO files..."
cp go.mod releases/arjunpreetham49-relay-datasource
cp Magefile.go releases/arjunpreetham49-relay-datasource
cp -r pkg releases/arjunpreetham49-relay-datasource

echo "Ziping up plugin..."
cd releases/
zip arjunpreetham49-relay-datasource-$PLUGIN_VERSION.zip . -r --exclude .DS_Store
rm -r arjunpreetham49-relay-datasource
zipinfo arjunpreetham49-relay-datasource-$PLUGIN_VERSION.zip
cd ../ 
rm -r dist
MD5_HASH=$(md5 releases/arjunpreetham49-relay-datasource-$PLUGIN_VERSION.zip)
echo "Done"

echo "MD5 Hash => $MD5_HASH"