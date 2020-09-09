import os
from setuptools import setup, find_packages

version = "0.1.11.dev0"

long_description = "\n\n".join([open("README.md").read(), open("CHANGES.md").read()])


def package_files(directory):
    paths = []
    for (path, directories, filenames) in os.walk(directory):
        for filename in filenames:
            paths.append(os.path.join("..", path, filename))
    return paths


# We use Pipenv. We set installation requirements here. Please set test requirements in Pipfile.
install_requirements = [
    "click-log==0.3.*",
    "click==7.*",
    "dacite==1.5.*",
    "dateparser==0.7.*",
    "lxml==4.*",
    "minio==4.0.*",
    "nested-dataclasses==0.1",
    "remotezip==0.9.*",
    "requests==2.22.*",
    "pystache==0.5.*",
]

style_template_files = package_files("atom_generator/resources/static_files")
xml_template_files = ["atom_generator/resources/templates/*"]
package_data = xml_template_files + style_template_files

setup(
    name="atom-generator",
    version=version,
    description="Atom Generator CLI to generate static (Inspire) Atom Feeds",
    long_description=long_description,
    # Get strings from http://www.python.org/pypi?%3Aaction=list_classifiers
    classifiers=["Programming Language :: Python :: 3"],
    keywords=["atom-generator"],
    author="Anton Bakker",
    author_email="anton.bakker@kadaster.nl",
    url="https://github.com/PDOK/atom-generator",
    packages=find_packages(exclude=["tests"]),
    package_data={"": package_data},
    setup_requires=["wheel"],
    include_package_data=True,
    zip_safe=False,
    install_requires=install_requirements,
    entry_points={
        "console_scripts": [
            "generate-atom = atom_generator.cli:generate_atom_service_command",
            "validate-models=atom_generator.cli:validate_models",
        ]
    },
)
