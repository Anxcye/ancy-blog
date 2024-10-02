package com.anxcye.handler.security;

import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.utils.WebUtils;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class AuthenticationEntryPointImpl implements AuthenticationEntryPoint {
    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response, AuthenticationException authException) throws IOException, ServletException {
        authException.printStackTrace();

        AppHttpCodeEnum result;
        if (authException instanceof InsufficientAuthenticationException) {
            result = AppHttpCodeEnum.NEED_LOGIN;
        } else if (authException instanceof BadCredentialsException) {
            result = AppHttpCodeEnum.LOGIN_ERROR;
        } else {
            result = AppHttpCodeEnum.AUTH_ERROR;
        }
        WebUtils.renderString(response, result);
    }
}
