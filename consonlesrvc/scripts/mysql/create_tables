#!/bin/bash

# crate tables defined in the directory "table"

HOSTNAME="localhost"
USERNAME="bcapp"
PASSWORD="DTXa9ytsNXpSEFuu6CPWc3ZxMcGiTpBq66mYyk9WYlI="

BATCH_FILE=`mktemp batchfile.XXXXXX`
echo "use app_wallet;" > $BATCH_FILE

for i in `ls tables`
do
  if [ -f tables/$i ]
  then
    echo "processing table $i"
    cat tables/$i >> $BATCH_FILE
  fi
done

echo "start creating tables..."
mysql -h$HOSTNAME -u$USERNAME -p$PASSWORD < $BATCH_FILE
echo "start creating tables...Done!"

echo "remove file $BATCH_FILE"
rm $BATCH_FILE
