package com.anxcye.domain.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * 网站设置
 * @TableName sys_setting
 */
@TableName(value ="sys_setting")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class Setting implements Serializable {
    /**
     * ID
     */
    @TableId(type = IdType.AUTO)
    private Long id;

    /**
     * 0正常 1停用
     */
    private String status;

    /**
     * 
     */
    private String name;

    /**
     * 
     */
    private Integer orderNum;

    /**
     * 
     */
    private String comment;

    /**
     * 1基础设置 2链接 3页脚
     */
    private Integer type;

    /**
     * 
     */
    private String value;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}