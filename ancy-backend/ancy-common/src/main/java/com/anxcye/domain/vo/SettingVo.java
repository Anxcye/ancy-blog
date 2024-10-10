package com.anxcye.domain.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class SettingVo {
    private String avatar;
    private String greeting;
    private String role;
    private String philosophy;
    private String name;
    private String address;
    private List<BadgeVo> badge;
    private List<FooterVo> footer;

    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class BadgeVo {
        private String index;
        private String title;
        private String url;
        private String img;
        private Integer orderNum;
    }

    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class FooterVo {
        private Integer position;
        private String index;
        private String text;
        private String url;
        private Integer orderNum;
    }

}
