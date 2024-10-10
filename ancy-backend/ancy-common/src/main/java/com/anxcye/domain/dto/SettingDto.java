package com.anxcye.domain.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class SettingDto {
    private String avatar;
    private String greeting;
    private String role;
    private String philosophy;
    private String name;
    private String address;
    private List<BadgeDto> badge;
    private List<FooterDto> footer;

    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class BadgeDto {
        private String index;
        private String title;
        private String url;
        private String img;
        private Integer orderNum;
    }

    @Data
    @NoArgsConstructor
    @AllArgsConstructor
    public static class FooterDto {
        private Integer position;
        private String index;
        private String text;
        private String url;
        private Integer orderNum;
    }

}
