import org.junit.jupiter.api.Test;

import com.anxcye.utils.JwtUtil;

public class TestJwt {

    @Test
    public void testJwt() {
        String jwt = JwtUtil.createJWT("admin");
        System.out.println(jwt);
    }
}
