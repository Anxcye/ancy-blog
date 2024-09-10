package com.anxcye.domain.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class categoryVo {
    private Long id;
    private String name;
    private Long parentId;
    private String description;
}
