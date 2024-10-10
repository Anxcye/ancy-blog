package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class CommentPageDto extends PageListDto implements Serializable {
    private Long articleId;
    private String email;
    private String nickname;
    private String content;
    private String status;
}