FROM python:3.8.7

COPY ./config.env .
COPY .docker/postgres_seed/ddl_seed.py . 
COPY ddl/ddl.sql .
COPY .docker/postgres_seed/requirements.txt . 

RUN pip3 install -r requirements.txt

CMD ["ddl_seed.py"]
ENTRYPOINT [ "python" ]