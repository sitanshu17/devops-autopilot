resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  availability_zone = "ap-south-1a"
  tags = {
    Name = "UbuntuInstance"
  }
}
