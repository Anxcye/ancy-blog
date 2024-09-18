package com.anxcye;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@MapperScan("com.anxcye.mapper")
public class AncyBlogApplication {

    public static void main(String[] args) {
        SpringApplication.run(AncyBlogApplication.class, args);
    }
}