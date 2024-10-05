# Mandatory Hand-in 2

Alice, Bob and Charlie are participating in a medical experiment in their local Hospital where they need to provide data for training a Machine Learning (ML) model, which will later be used to diagnose patients.
In this setting, the Hospital runs a central server that collects data from all patients.
The Hospital is trusted to honestly execute a data collection protocol pre-agreed among all patients taking part in the experiment.
However, the Hospital is aware of security issues and does not want to store this data in its systems.
Moreover, do not trust each other to see their data, but do trust each other to follow the protocol established by the Hospital.
As usual, the patients and the Hospital communicate over the Internet.

The restrictions in the scenario above leave the Hospital and patients in a tight spot, since standard ML algorithms require access to plaintext data in order to train a model.
Luckily, the researchers in the hospital are collaborating with a team of security experts who suggested using a Federated Learning algorithm that supports Secure Aggregation, which allows the algorithm to train a model while seeing only an aggregate value computed from all patients' data.
In particular, this algorithm works on aggregated values obtained by summing all individual patient's values.
Using this technique, neither the patients nor the Hospital get access to patient's individual plaintext data, but only to aggregated values.

1) Reflect on this scenario in the context of the GDPR:
    What are the potential issues in having the hospital store plaintext private data provided by patients even if they have consented to participate on the experiment and have their data processed?
    Would these issues be solved by removing the patients' names from their data before storing it?
    What are the remaining risks in using Federated Learning with Secure Aggregation as suggested?

2) Design and implement a solution that allows for the 3 patients and the Hospital in the scenario above to compute an aggregate value of the patients' private input values in such a way that the Hospital only learns this aggregate value and no patient learns anything besides their own private inputs.
Your protocol must also ensure confidentiality and integrity of the data against external attackers.
Consider that all individual values held by patients are integers in a range [0,...,R] and that the aggregated value is the sum of all individual values, which is also assumed to be in the same range.
You must describe an adversarial model (or threat model) capturing the attacks by an adversary who behaves according to the scenario described above, explain how your solution works and why it guarantees security against such an adversary.

Hint: Secure Aggregation for Federated Learning is a real-world practical technique.

Deliverables:

- A written report reflecting on the questions in "1", describing the adversarial model you are working on, describing the building blocks of your proposed solution, how they are combined in your final solution and why they guarantee security against the adversary you describe.
- An implementation of your solution in a programming language of your choosing, along with clear instructions on how to compile and run it (potentially added to the report or to a separate Readme file).

Submission Instructions:

- Submit only the PDF file with your report and the file containing your code. Do not submit whole folders containing metadata, auxiliary IDE files or anything else than the code and report.
- Please name your submission clearly using your Name/Student ID, e.g. Jane Doe - 36476832.zip, Jane Doe - 36476832 - Report.PDF, Jane Doe - 36476832 - code.c, Jane Doe - 36476832 - Readme.txt. This makes grading faster, so that you get feedback on your hand-in faster.
