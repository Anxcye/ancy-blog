package com.anxcye.service;

import com.anxcye.domain.dto.SettingDto;
import com.anxcye.domain.entity.Setting;
import com.anxcye.domain.vo.SettingVo;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【sys_setting(网站设置)】的数据库操作Service
* @createDate 2024-10-04 13:12:16
*/
public interface SettingService extends IService<Setting> {

    SettingVo getBaseSetting();

    boolean updateBaseSetting(SettingDto settingDto);
}
