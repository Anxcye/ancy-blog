package com.anxcye.aspect;

import com.alibaba.fastjson.JSON;
import com.anxcye.annotation.Log;
import com.anxcye.domain.entity.OperateLog;
import com.anxcye.service.OperateLogService;
import com.anxcye.utils.SecurityUtil;
import jakarta.servlet.http.HttpServletRequest;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Pointcut;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.lang.reflect.Method;
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

    private String getCustomFieldValue(ProceedingJoinPoint joinPoint) {
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        Method method = signature.getMethod();
        Log logAnnotation = method.getAnnotation(Log.class);
        
        if (logAnnotation != null && !logAnnotation.fieldName().isEmpty()) {
            String fieldName = logAnnotation.fieldName();
            Object[] args = joinPoint.getArgs();
            if (args.length > 0 && args[0] != null) {
                Object arg = args[0];  // 获取第一个参数
                try {
                    java.lang.reflect.Field field = arg.getClass().getDeclaredField(fieldName);
                    field.setAccessible(true);
                    Object value = field.get(arg);
                    if (value != null) {
                        return value.toString();
                    }
                } catch (NoSuchFieldException | IllegalAccessException e) {
                    // 字段不存在或无法访问
                }
            }
        }
        return "";
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

        String params = Arrays.toString(joinPoint.getArgs());

        String methodParams = params.length() > 1000 ? params.substring(0, 995) + "..." : params;

        String summaryValue = getCustomFieldValue(joinPoint);

        String returnV = JSON.toJSONString(result);

        String returnValue = returnV.length() > 2000 ? returnV.substring(0, 1995) + "..." : returnV;

        long costTime = endTime - startTime;

        OperateLog operateLog = new OperateLog(
                null,
                operateUser,
                ip,
                address,
                ua,
                operateTime,
                summaryValue,
                className,
                methodName,
                methodParams,
                returnValue,
                costTime);

        operateLogService.save(operateLog);

        return result;
    }

}
