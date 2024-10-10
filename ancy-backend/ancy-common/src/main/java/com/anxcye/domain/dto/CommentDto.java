package com.anxcye.domain.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class CommentDto {

    /**
     * 评论类型（0代表文章评论，1代表NOTE评论）
     */
    private String type;

    /**
     * 文章id
     */
    private Long articleId;

    /**
     * 根评论id
     */
    private Long parentId;

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
     * 公开状态 0代表公开，1代表隐藏
     */

    private String status;
    /**
     *
     */
    private Integer likeCount;

    /**
     *
     */
    private String isTop;

    /**
     * 所回复的目标评论的userid
     */
    private String toCommentNickname;

    /**
     * 回复目标评论id
     */
    private Long toCommentId;


}
