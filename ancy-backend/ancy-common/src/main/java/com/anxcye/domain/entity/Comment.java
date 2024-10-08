package com.anxcye.domain.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;

import java.io.Serializable;
import java.util.Date;

/**
 * 评论表
 * @TableName ancy_comment
 */
@TableName(value ="ancy_comment")
@Data
public class Comment implements Serializable {
    /**
     * 
     */
    @TableId(type = IdType.AUTO)
    private Long id;

    /**
     * 评论类型（0代表文章评论，1代表NOTE评论）
     */
    private String type;

    /**
     * 文章id
     */
    private Long articleId;

    /**
     * 公开状态 0代表公开，1代表隐藏
     */
    private String status;

    /**
     * 根评论id
     */
    private Long parentId;

    /**
     * 
     */
    private Long userId;

    /**
     * 
     */
    private String avatar;

    /**
     * 
     */
    private String nickname;

    /**
     * 
     */
    private String email;

    /**
     * 评论内容
     */
    private String content;

    /**
     * 
     */
    private String ua;

    /**
     * 
     */
    private String ip;

    /**
     * 
     */
    private Integer likeCount;

    /**
     * 
     */
    private String isTop;

    /**
     * 
     */
    private String toCommentNickname;

    /**
     * 
     */
    private Long toCommentId;

    /**
     * 
     */
    @TableField(fill = FieldFill.INSERT)
    private Long createBy;

    /**
     * 
     */
    @TableField(fill = FieldFill.INSERT)
    private Date createTime;

    /**
     * 
     */
    @TableField(fill = FieldFill.UPDATE)
    private Long updateBy;

    /**
     * 
     */
    @TableField(fill = FieldFill.UPDATE)
    private Date updateTime;

    /**
     * 删除标志（0代表未删除，1代表已删除）
     */
    private Integer deleted;

    @TableField(exist = false)
    private static final long serialVersionUID = 1L;
}