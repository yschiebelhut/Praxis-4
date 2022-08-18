{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action":  "s3:*",
            "Resource": [
                "arn:aws:s3:::ai-terraform-{{.}}",
                "arn:aws:s3:::ai-terraform-{{.}}/*",
                "arn:aws:s3:::backup-cluster-{{.}}",
                "arn:aws:s3:::backup-cluster-{{.}}/*",
                "arn:aws:s3:::mlf-gitops-artifact-registry-{{.}}",
                "arn:aws:s3:::mlf-gitops-artifact-registry-{{.}}/*",
                "arn:aws:s3:::mlf-gitops-artifacts-registry-{{.}}",
                "arn:aws:s3:::mlf-gitops-artifacts-registry-{{.}}/*",
                "arn:aws:s3:::bucket-logs-loki-{{.}}",
                "arn:aws:s3:::bucket-logs-loki-{{.}}/*",
                "arn:aws:s3:::aicore-acceptance-tests-{{.}}",
                "arn:aws:s3:::aicore-acceptance-tests-{{.}}/*",
                "arn:aws:s3:::aicore-nvidia-driver-{{.}}",
                "arn:aws:s3:::aicore-nvidia-driver-{{.}}/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action":  "iam:*",
            "Resource": [
                "arn:aws:iam:::user/ai-gitops-velero-{{.}}--bot",
                "arn:aws:iam:::user/mlf-gitops-artifact-registry-{{.}}--bot",
                "arn:aws:iam:::policy/mlf-gitops-artifact-registry-{{.}}--policy",
                "arn:aws:iam:::policy/ai-gitops-velero-{{.}}--policy",
                "arn:aws:iam:::user/aicore-acceptance-tests-{{.}}--bot",
                "arn:aws:iam:::policy/aicore-acceptance-tests-{{.}}--policy",
                "arn:aws:iam:::user/aicore-nvidia-driver-read-{{.}}--bot",
                "arn:aws:iam:::policy/aicore-nvidia-driver-read-{{.}}--policy",
                "arn:aws:iam:::user/ai-gitops-route53-{{.}}--bot",
                "arn:aws:iam:::policy/ai-gitops-route53-{{.}}--policy",
                "arn:aws:iam:::user/ai-gitops-route53-certmanager-{{.}}-bot",
                "arn:aws:iam:::policy/ai-gitops-route53-certmanager-{{.}}-policy",
                "arn:aws:iam:::user/ai-gitops-loki-{{.}}--bot",
                "arn:aws:iam:::policy/ai-gitops-loki-{{.}}--policy"
            ]
        }
    ]
}