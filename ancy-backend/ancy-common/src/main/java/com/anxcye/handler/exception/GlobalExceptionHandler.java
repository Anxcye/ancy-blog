package com.anxcye.handler.exception;

import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.exception.SystemException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
@Slf4j
public class GlobalExceptionHandler {

    @ExceptionHandler(SystemException.class)
    public ResponseResult<?> systemExceptionHandler(SystemException e) {
        log.error("SystemException: ", e.getAppHttpCodeEnum().getMsg());
        return ResponseResult.error(e.getAppHttpCodeEnum());
    }

    @ExceptionHandler(RuntimeException.class)
    public ResponseResult<?> runtimeExceptionHandler(Exception e) throws Exception {
        if (e instanceof BadCredentialsException ||
                e instanceof AccessDeniedException
        ) {
            throw e;
        }
        log.error("SystemException: ", e);
        return ResponseResult.error(AppHttpCodeEnum.SYSTEM_ERROR.getCode(), e.getMessage());
    }
}
