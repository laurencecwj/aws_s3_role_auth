from minio import Minio
from minio.credentials import IamAwsProvider
import configparser
import io

config = configparser.ConfigParser()
config.read('./config.ini')

bucket_name = config['aws']['s3_bucket']
print("using bucket: ", bucket_name)
s3_endpoint = config['aws'].get('s3_endpoint', 's3.cn-northwest-1.amazonaws.com.cn')
print("using s3 endpoint: ", s3_endpoint)
object_name = "py.foo.bar"
print("using s3 object: ", object_name)

provider = IamAwsProvider()
client = Minio(
    s3_endpoint,
    credentials=provider,
    secure=True
)

def put_object():
    try:
        data = b"Hello, MinIO!"
        data_stream = io.BytesIO(data)
    
        client.put_object(
            bucket_name,
            object_name,
            data_stream,
            length=len(data),
            content_type="text/plain"
        )
        print(f"'{object_name}' uploaded successfully to '{bucket_name}'.")
    except Exception as e:
        print(f"put object to {bucket_name}/{object_name} error: {e}")


def get_object():
    try:
        response = client.get_object(bucket_name, object_name)
        object_data = response.read()
        print(f"Downloaded '{object_name}': {object_data.decode('utf-8')}")
    except Exception as e:
        print(f"Error getting {bucket_name}/{object_name} object error: {e}")
    finally:
        response.close()
        response.release_conn()

def rm_object():
    try:
        client.remove_object(bucket_name, object_name)
        print(f"'{object_name}' deleted successfully from '{bucket_name}'.")
    except Exception as e:
        print(f"Error deleting {bucket_name}/{object_name} object: {e}")

put_object()
get_object()
rm_object()
