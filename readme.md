

### Setup
1. `docker build compose up`

### Running 
1. CLI, two input methods
   1. stdin example: `cat resources/query_params.csv | docker-compose run --rm cli`
   2. command line argument example: `docker-compose run cli --rm  -f resources/query_params_no_header.csv` 
   
