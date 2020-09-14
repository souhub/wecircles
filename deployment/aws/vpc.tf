# _____________________________VPC_______________________________________
resource "aws_vpc" "vpc" {
  cidr_block = "172.0.0.0/16"
  tags = {
    Name = "wecircles"
  }
}


# _____________________________Subnet_______________________________________
resource "aws_subnet" "publicA" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.0.1.0/24"
  availability_zone = "ap-northeast-1a"
  tags = {
    Name = "publicA"
  }
}

resource "aws_subnet" "publicC" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.0.2.0/24"
  availability_zone = "ap-northeast-1c"
  tags = {
    Name = "publicC"
  }
}

resource "aws_subnet" "privateA" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.0.10.0/24"
  availability_zone = "ap-northeast-1a"
  tags = {
    Name = "privateA"
  }
}

resource "aws_subnet" "privateC" {
  vpc_id            = aws_vpc.vpc.id
  cidr_block        = "172.0.20.0/24"
  availability_zone = "ap-northeast-1c"
  tags = {
    Name = "privateC"
  }
}

# _____________________________Internet Gateway_______________________________________
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "wecircles"
  }
}

# _____________________________Route Table_______________________________________
resource "aws_route_table" "rt" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "wecircles"
  }
}

# Internet Gateway Association
resource "aws_route" "public" {
  route_table_id         = aws_route_table.rt.id
  gateway_id             = aws_internet_gateway.igw.id
  destination_cidr_block = "0.0.0.0/0"
}

# Subnet Associations
resource "aws_route_table_association" "public_0" {
  subnet_id      = aws_subnet.publicA.id
  route_table_id = aws_route_table.rt.id
}

resource "aws_route_table_association" "public_1" {
  subnet_id      = aws_subnet.publicC.id
  route_table_id = aws_route_table.rt.id
}
