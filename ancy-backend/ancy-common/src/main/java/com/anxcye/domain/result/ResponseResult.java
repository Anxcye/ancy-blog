package com.anxcye.domain.result;

import com.anxcye.domain.enums.AppHttpCodeEnum;
import lombok.Data;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

import java.io.Serializable;

@Data
public class ResponseResult<T> implements Serializable {

    private Integer code;
    private String msg;
    private T data;


    public static <T> ResponseResult<T> success(AppHttpCodeEnum appHttpCodeEnum, T data) {
        ResponseResult<T> result = new ResponseResult<>();
        result.code = appHttpCodeEnum.getCode();
        result.msg = appHttpCodeEnum.getMsg();
        result.data = data;
        return result;
    }

    public static <T> ResponseResult<T> success() {
        return success(AppHttpCodeEnum.SUCCESS, null);
    }

    public static <T> ResponseResult<T> success(T object) {
        return success(AppHttpCodeEnum.SUCCESS, object);
    }

    public static <T> ResponseResult<T> error(AppHttpCodeEnum appHttpCodeEnum) {
        return error(appHttpCodeEnum.getCode(), appHttpCodeEnum.getMsg());
    }

    public static <T> ResponseResult<T> error(Integer code, String msg) {
        ResponseResult<T> result = new ResponseResult<>();
        result.msg = msg;
        result.code = code;
        ServletRequestAttributes attributes = ((ServletRequestAttributes) RequestContextHolder.getRequestAttributes());
        if (attributes != null) {
            assert attributes.getResponse() != null;
            attributes.getResponse().setStatus(code);
        }
        return result;
    }

}
