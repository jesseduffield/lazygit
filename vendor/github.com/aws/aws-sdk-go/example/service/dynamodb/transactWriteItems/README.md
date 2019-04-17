# AWS DynamoDB Transaction Error Aware Client for Go
The client provides a workaround for [this bug](https://github.com/aws/aws-sdk-go/issues/2318)

## How to use
This example shows how to use the client to read transaction error cancellation reasons.
```go
    sess := session.Must(session.NewSession())
    svc := NewTxErrorAwareDynamoDBClient(sess)

    input := &dynamodb.TransactWriteItemsInput{
    	//...
    }

    if _, err := svc.TransactWriteItems(input); err != nil {
    	txErr := err.(TxRequestFailure)
    	fmt.Println(txErr.CancellationReasons())
    }
```

Sample response of the Println statement
```
{com.amazonaws.dynamodb.v20120810#TransactionCanceledException Transaction cancelled, please refer cancellation reasons for specific reasons [ConditionalCheckFailed, None, None] [{
  Code: "ConditionalCheckFailed",
  Item: {
    AlbumTitle: {
      S: "==========     43"
    },
    Artist: {
      S: "Acme Band 14"
    },
    Year: {
      N: "2017"
    },
    SongTitle: {
      S: "Happy Day 12"
    }
  },
  Message: "The conditional request failed"
} {
  Code: "None"
} {
  Code: "None"
}]}
[{
  Code: "ConditionalCheckFailed",
  Item: {
    AlbumTitle: {
      S: "==========     43"
    },
    Artist: {
      S: "Acme Band 14"
    },
    Year: {
      N: "2017"
    },
    SongTitle: {
      S: "Happy Day 12"
    }
  },
  Message: "The conditional request failed"
} {
  Code: "None"
} {
  Code: "None"
}]
```
