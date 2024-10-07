package com.anxcye.domain.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.io.Serializable;
import java.util.Date;
import lombok.Data;

/**
 * 阅读表
 * @TableName ancy_read
 */
@TableName(value ="ancy_read")
@Data
public class Read implements Serializable {
    /**
     * 
     */
    @TableId(type = IdType.AUTO)
    private Long id;

    /**
     * 出处
     */
    private String source;

    /**
     * 内容
     */
    private String content;

    /**
     * 作者
     */
    private String author;

    /**
     * 添加来源 0 手动添加 1安读
     */
    private Integer addFrom;

    /**
     * 
     */
    private Long createBy;

    /**
     * 
     */
    private Date createTime;

    /**
     * 
     */
    private Long updateBy;

    /**
     * 
     */
    private Date updateTime;

    /**
     * 删除标志（0代表未删除，1代表已删除）
     */
    private Integer deleted;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}