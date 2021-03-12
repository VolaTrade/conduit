FROM python:3.8.7

COPY ddl/seed.py . 
COPY ddl/ddl.sql .
COPY ddl/requirements.txt .
COPY config.env . 

RUN pip3 install -r requirements.txt

CMD ["seed.py"]
ENTRYPOINT [ "python" ]