# AWS Discover Test Tool

A simple CLI tool to test the AWS discover provider, particularly the dual-stack endpoint fix.

## Build

### Local (macOS)

```bash
cd cmd/test-aws-discover
go build -o test-aws-discover .
```

### Cross-compile for Linux amd64 (to run on AWS EC2)

```bash
cd cmd/test-aws-discover
GOOS=linux GOARCH=amd64 go build -o test-aws-discover-linux-amd64 .
```

## Running on AWS EC2

Copy the binary to your EC2 instance and run it. The instance's IAM role will be used for authentication.

```bash
# From your Mac, copy to EC2
scp test-aws-discover-linux-amd64 ec2-user@<instance-ip>:~/

# SSH to EC2 and run
ssh ec2-user@<instance-ip>
chmod +x ~/test-aws-discover-linux-amd64

# Test with dual-stack disabled (the fix)
AWS_USE_DUALSTACK_ENDPOINT=false ~/test-aws-discover-linux-amd64 -region me-central-1 -tag-key Name -tag-value my-instance

# Test with dual-stack enabled (default behavior)
~/test-aws-discover-linux-amd64 -region us-east-1 -tag-key Name -tag-value my-instance
```

## Usage

```bash
./test-aws-discover -region <region> -tag-key <key> -tag-value <value> [-addr-type <type>]
```

## Test Scenarios

### 1. Test in us-east-1 with dual-stack enabled (default)

```bash
./test-aws-discover -region us-east-1 -tag-key Name -tag-value my-instance
```

### 2. Test in me-central-1 with dual-stack disabled (THE FIX)

```bash
# This should work now with the fix
AWS_USE_DUALSTACK_ENDPOINT=false ./test-aws-discover -region me-central-1 -tag-key Name -tag-value my-instance
```

### 3. Test in me-central-1 with dual-stack enabled (will fail - no dual-stack in this region)

```bash
# This will fail because me-central-1 doesn't have dual-stack endpoints
./test-aws-discover -region me-central-1 -tag-key Name -tag-value my-instance
```

## Expected Behavior

| Region | `AWS_USE_DUALSTACK_ENDPOINT` | Expected Result |
|--------|------------------------------|-----------------|
| us-east-1 | `true` or unset | Works (dual-stack available) |
| us-east-1 | `false` | Works (standard endpoints) |
| me-central-1 | `false` | **Works (THE FIX)** |
| me-central-1 | `true` or unset | Fails (no dual-stack in this region) |
