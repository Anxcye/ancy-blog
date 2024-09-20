package com.anxcye;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@MapperScan("com.anxcye.mapper")
@EnableScheduling
public class AncyAdminApplication {

    public static void main(String[] args) {
        SpringApplication.run(AncyAdminApplication.class, args);
    }
}