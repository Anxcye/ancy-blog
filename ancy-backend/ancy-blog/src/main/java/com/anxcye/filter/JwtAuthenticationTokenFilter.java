package com.anxcye.filter;

import com.anxcye.utils.RedisCache;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;

@Component
public class JwtAuthenticationTokenFilter extends OncePerRequestFilter {

    @Autowired
    private RedisCache redisCache;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {

//
//        String token = request.getHeader("token");
//        if (StringUtils.isEmpty(token)) {
//            filterChain.doFilter(request, response);
//            return;
//        }
//        LoginUser loginUser;
//        try {
//            Claims claims = JwtUtil.parseJWT(token);
//            String userId = claims.getSubject();
//            JSONObject jsonObject = redisCache.getCacheObject(RedisConstant.BLOG_TOKEN_PREFIX + userId);
//            loginUser = jsonObject.toJavaObject(LoginUser.class);
//        } catch (Exception e) {
//            WebUtils.renderString(response, AppHttpCodeEnum.TOKEN_INVALID);
//            return;
//        }
//
//        UsernamePasswordAuthenticationToken authenticationToken =
//                new UsernamePasswordAuthenticationToken(loginUser, null, null);
//        SecurityContextHolder.getContext().setAuthentication(authenticationToken);

        filterChain.doFilter(request, response);
    }
}
