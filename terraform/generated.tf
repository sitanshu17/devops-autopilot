```terraform
resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  region        = "ap-south-1"
  tags = {
    Name = "UbuntuInstance"
  }
}
```