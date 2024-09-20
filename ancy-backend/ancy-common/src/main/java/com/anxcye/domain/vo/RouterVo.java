package com.anxcye.domain.vo;

import com.anxcye.domain.vo.MenuVo.MenuVo;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class RouterVo {
    private List<MenuVo> menus;
}
