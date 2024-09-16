import com.anxcye.utils.JwtUtil;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jws;
import io.jsonwebtoken.Jwts;
import org.junit.jupiter.api.Test;

import javax.crypto.SecretKey;

public class TestJwt {

    @Test
    public void testJwt() {
        String jwt = JwtUtil.createJWT("admin");
        System.out.println(jwt);
    }

    @Test
    public void testJwt2() {
        Claims claim = null;
        try {
            claim = JwtUtil.parseJWT("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJqb2UifQ.1lnR-gjeq1S5G5PklO4mVDo7hcd5ovQ-BB-MT68QgP4");
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println(claim);
    }

    @Test
    public void createAndParse() throws Exception {
        String  jwt = JwtUtil.createJWT("admin");
        Claims claims = JwtUtil.parseJWT(jwt);
        System.out.println(claims);

    }

    @Test
    public void jwtsTest(){
        SecretKey key = Jwts.SIG.HS256.key().build();
        String jws = Jwts.builder().subject("Joe").signWith(key).compact();

        Jws<Claims> claimsJws = Jwts.parser().verifyWith(key).build().parseSignedClaims(jws);

        System.out.println(claimsJws);


    }
}
