# no-phi-ai: Advanced PHI/PII Detection in GitHub Activities

## Introduction

no-phi-ai is a GitHub App developed in GoLang, designed to enhance the security of GitHub repositories by detecting and preventing the exposure of Personal Health Information (PHI) and Personally Identifiable Information (PII). Leveraging the power of Azure AI Language service, it meticulously scans GitHub comments, issues, and pull requests, ensuring sensitive data is identified and managed appropriately.

## Azure AI Language Service Integration

The integration with Azure AI Language service significantly boosts the app's ability to detect PHI/PII within textual content. By utilizing advanced natural language processing technologies, the service analyzes text for sensitive information, providing a robust layer of protection against potential data breaches.

## Core Capabilities

no-phi-ai offers a comprehensive suite of features aimed at actively monitoring and protecting GitHub repositories:
- **Comments Analysis:** Scans comments across commits, pull requests, and issues for PHI/PII.
- **Issues Monitoring:** Reviews issue content for sensitive information.
- **Pull Requests Inspection:** Examines pull request changes to detect PHI/PII.

## Mitigating PHI/PII Exposure

Upon detecting PHI/PII, no-phi-ai takes proactive steps to mitigate exposure risks:
- **Alerts:** Notifies repository administrators about potential PHI/PII exposure.
- **Redaction Options:** Provides mechanisms to redact sensitive information from the content.
- **Secure Handling Recommendations:** Offers advice on managing detected information securely.

## Installation and Configuration

To integrate no-phi-ai with your GitHub repository:
1. Visit the GitHub Marketplace and search for "no-phi-ai".
2. Follow the installation process and select the repositories you wish to protect.
3. Customize the scanning settings via the app's dashboard to suit your project's requirements.

## Contributing to no-phi-ai

We welcome contributions to no-phi-ai, which play a vital role in improving its functionality and reach. Whether it's fixing bugs, adding new features, or enhancing documentation, your contributions are greatly appreciated. Please refer to our contributing guidelines for more information on how to contribute.

## Additional Resources

For further information on no-phi-ai and how to protect your projects against PHI/PII exposure, explore the following resources:
- [Project License](LICENSE)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Azure AI Language Service Documentation](https://docs.microsoft.com/en-us/azure/cognitive-services/language-service/)

Inspired by the GitHub Action from rob-derosa/pii-detection, no-phi-ai aims to provide an enhanced security layer for GitHub repositories, ensuring the confidentiality and integrity of sensitive data.
