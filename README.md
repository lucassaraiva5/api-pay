# api-pay

## Prerequisites

Before running the project, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)


## Clone the project

First need to clone the project with the command:

```
cd ~ && git clone https://github.com/lucassaraiva5/api-pay
```

## Navigate to the folder

Now you need to navigate to the project folder:

```
cd ~/api-pay
```

## Getting Started

Begin by duplicating the `.env.example` file and renaming the copy to `.env`. This will set up your environment variables for local development.

```
cp .env.example .env
```

## Run the project with docker

To quickly start the project, simply run the following command:

```
docker compose up -d
```

This will build and launch all necessary services in detached mode.

## How to test it

You can use the provided Insomnia collection to test the API endpoints. The collection file is located in the project directory. Import this file into Insomnia to access pre-configured requests for the API.

For more information on importing collections, see the [Insomnia documentation](https://docs.insomnia.rest/insomnia/import-export-data).