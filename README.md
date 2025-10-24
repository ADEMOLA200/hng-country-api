# HNG Country API

A RESTful API that fetches country data from external APIs, stores it in a database, and provides CRUD operations with GDP calculations.

## Features

- Fetch country data from REST Countries API
- Get exchange rates from Open Exchange Rates API
- Calculate estimated GDP based on population and exchange rates
- CRUD operations for country data
- Filtering and sorting capabilities
- Automatic image generation with summary statistics
- MySQL database integration

## Setup Instructions

### Prerequisites

- Go 1.21+
- MySQL 5.7+

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd hng-country-api
