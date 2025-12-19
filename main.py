from minio import Minio
from minio.credentials import IamAwsProvider
import configparser
import io
import boto3

config = configparser.ConfigParser()
config.read('./config.ini')

bucket_name = config['aws']['s3_bucket']
print("using bucket: ", bucket_name)
s3_endpoint = config['aws'].get('s3_endpoint', 's3.cn-northwest-1.amazonaws.com.cn')
print("using s3 endpoint: ", s3_endpoint)
object_name = "py.foo.bar"
print("using s3 object: ", object_name)
region = config['aws'].get('s3_region', 'cn-northwest-1')
print("using s3 region: ", region)

provider = IamAwsProvider()
client = Minio(
    s3_endpoint,
    credentials=provider,
    secure=True
)

client2 = boto3.resource('s3', region_name=region).Bucket(bucket_name)

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

def put_object_boto3():
    try:
        data = b"Hello, MinIO!"
        data_stream = io.BytesIO(data)
    
        client2.put_object(
            Body=data_stream, Bucket=bucket_name, Key=object_name
        )
        print(f"'{object_name}' uploaded successfully to '{bucket_name}'.")
    except Exception as e:
        print(f"put object to {bucket_name}/{object_name} error: {e}")


def get_object_boto3():
    try:
        obj = client2.Object(bucket_name=bucket_name, key=object_name)
        object_data = obj.get()["Body"].read()
        print(f"Downloaded '{object_name}': {object_data.decode('utf-8')}")
    except Exception as e:
        print(f"Error getting {bucket_name}/{object_name} object error: {e}")
    finally:
        pass

def rm_object_boto3():
    try:
        client2.Object(object_name).delete()
        print(f"'{object_name}' deleted successfully from '{bucket_name}'.")
    except Exception as e:
        print(f"Error deleting {bucket_name}/{object_name} object: {e}")

print("testing minio sdk")
put_object()
get_object()
rm_object()

print("testing boto3 sdk")
put_object_boto3()
get_object_boto3()
rm_object_boto3()
