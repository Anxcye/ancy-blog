package com.anxcye;

import com.anxcye.utils.JwtUtil;
import io.jsonwebtoken.Claims;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
public class TestJwt {

    @Test
    public void testJwt() {
        String jwt = JwtUtil.createJWT("admin");
        System.out.println(jwt);
    }

//    @Test
//    public void testJwt2() {
//        Claims claim = null;
//        try {
//            claim = JwtUtil.parseJWT("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJqb2UifQ.1lnR-gjeq1S5G5PklO4mVDo7hcd5ovQ-BB-MT68QgP4");
//        } catch (Exception e) {
//            e.printStackTrace();
//        }
//        System.out.println(claim);
//    }

    @Test
    public void createAndParse() throws Exception {
        String  jwt = JwtUtil.createJWT("admin");
        Claims claims = JwtUtil.parseJWT(jwt);
        System.out.println(claims);

    }

}
