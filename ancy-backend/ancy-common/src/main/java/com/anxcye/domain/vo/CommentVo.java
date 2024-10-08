package com.anxcye.domain.vo;

import com.anxcye.domain.result.PageResult;
import lombok.Data;

import java.io.Serializable;
import java.util.Date;

@Data
public class CommentVo implements Serializable {
    private Long id;

    /**
     * 文章id
     */
    private Long articleId;

    /**
     * 审核状态 (0代表审核通过，1代表审核未通过，2代表未审核)
     */
    private String status;

    /**
     * 根评论id
     */
    private Long parentId;

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
    private PageResult children;

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
}