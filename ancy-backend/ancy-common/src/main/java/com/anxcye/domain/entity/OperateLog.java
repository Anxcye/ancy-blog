package com.anxcye.domain.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import java.util.Date;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 操作日志表
 * @TableName sys_operate_log
 */
@TableName(value ="sys_operate_log")
@Data
@AllArgsConstructor
@NoArgsConstructor
public class OperateLog implements Serializable {
    /**
     * ID
     */
    @TableId(type = IdType.AUTO)
    private Long id;

    /**
     * 操作人ID
     */
    private Long operateUser;

    /**
     * 操作人IP
     */
    private String ip;

    /**
     * 操作人地址
     */
    private String address;

    /**
     * 操作人UA
     */
    private String ua;

    /**
     * 操作时间
     */
    private Date operateTime;

    /**
     * 操作的类名
     */
    private String className;

    /**
     * 操作的方法名
     */
    private String methodName;

    /**
     * 方法参数
     */
    private String methodParams;

    /**
     * 返回值
     */
    private String returnValue;

    /**
     * 方法执行耗时, 单位:ms
     */
    private Long costTime;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}