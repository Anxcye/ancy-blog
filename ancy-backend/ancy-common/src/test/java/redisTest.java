import com.anxcye.constants.RedisConstant;
import com.anxcye.domain.entity.LoginUser;
import com.anxcye.utils.RedisCache;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.RedisTemplate;

@ExtendWith(MockitoExtension.class)
public class redisTest {

    @Mock
    private RedisTemplate redisTemplate;

    @InjectMocks
    private RedisCache redisCache;

    @Test
    public void redisTest() {
        LoginUser loginUser = redisCache.getCacheObject(RedisConstant.BLOG_TOKEN_PREFIX + 1);
        System.out.println(loginUser.getUser().getId());

    }
}
