package com.anxcye.aspect;

import com.alibaba.fastjson.JSON;
import com.anxcye.domain.entity.OperateLog;
import com.anxcye.service.OperateLogService;
import com.anxcye.utils.SecurityUtil;
import jakarta.servlet.http.HttpServletRequest;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Pointcut;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.Arrays;
import java.util.Date;

@Aspect
@Component
public class LogAspect {

    @Autowired
    private OperateLogService operateLogService;

    @Autowired
    private HttpServletRequest request;

    @Pointcut("@annotation(com.anxcye.annotation.Log)")
    public void logPointCut() {
    }

    @Around("logPointCut()")
    public Object around(ProceedingJoinPoint joinPoint) throws Throwable {

        long startTime = System.currentTimeMillis();
        Object result = joinPoint.proceed();
        long endTime = System.currentTimeMillis();
        Long operateUser;
        try {
            operateUser = SecurityUtil.getUserId();
        } catch (Exception e) {
            operateUser = -1L;
        }

        String ip = request.getRemoteAddr();

        String address = request.getRemoteAddr();

        String ua = request.getHeader("User-Agent");

        Date operateTime = new Date();

        String className = joinPoint.getTarget().getClass().getName();

        String methodName = joinPoint.getSignature().getName();

        String methodParams = Arrays.toString(joinPoint.getArgs());

        String returnValue = JSON.toJSONString(result);

        long costTime = endTime - startTime;

        OperateLog operateLog = new OperateLog(
                null,
                operateUser,
                ip,
                address,
                ua,
                operateTime,
                className,
                methodName,
                methodParams,
                returnValue,
                costTime);

        operateLogService.save(operateLog);

        return result;
    }

}
