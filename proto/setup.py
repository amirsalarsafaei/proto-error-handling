from setuptools import setup, find_packages

setup(
    name="proto-error-handling-interface",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "protobuf>=4.21.0",
    ],
    author="Your Name",
    author_email="your.email@example.com",
    description="Generated Protocol Buffers for Python",
    url="https://github.com/amirsalarsafaei/proto-error-handling",
    classifiers=[
        "Programming Language :: Python :: 3",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
)
