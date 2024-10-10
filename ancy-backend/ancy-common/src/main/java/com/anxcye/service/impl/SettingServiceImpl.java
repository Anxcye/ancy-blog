package com.anxcye.service.impl;

import com.alibaba.fastjson.JSONObject;
import com.alibaba.fastjson.TypeReference;
import com.anxcye.annotation.Log;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.SettingDto;
import com.anxcye.domain.entity.Setting;
import com.anxcye.domain.vo.SettingVo;
import com.anxcye.mapper.SettingMapper;
import com.anxcye.service.SettingService;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * @author axy
 * @description 针对表【sys_setting(网站设置)】的数据库操作Service实现
 * @createDate 2024-10-04 13:12:16
 */
@Service
public class SettingServiceImpl extends ServiceImpl<SettingMapper, Setting>
        implements SettingService {

    private LambdaQueryWrapper<Setting> getWrapper() {
        LambdaQueryWrapper<Setting> wrapper = new LambdaQueryWrapper<>();
        if (!SecurityUtil.isAdmin()) {
            wrapper.eq(Setting::getStatus, SystemConstants.STATUS_NORMAL);
        }
        return wrapper;
    }

    @Override
    public SettingVo getBaseSetting() {
        LambdaQueryWrapper<Setting> wrapper = getWrapper();
        wrapper.orderByAsc(Setting::getOrderNum);
        List<Setting> settingList = list(wrapper);

        Map<String, String> infoMap = new HashMap<>();
        List<SettingVo.BadgeVo> badgeList = new ArrayList<>();
        List<SettingVo.FooterVo> footerList = new ArrayList<>();

        settingList.forEach(setting -> {
            switch (setting.getType()) {
                case 1:
                    infoMap.put(setting.getName(), setting.getValue());
                    break;
                case 2:
                    JSONObject badge = JSONObject.parseObject(setting.getValue());
                    SettingVo.BadgeVo badgeVo = badge.toJavaObject(SettingVo.BadgeVo.class);
                    badgeVo.setOrderNum(setting.getOrderNum());
                    badgeList.add(badgeVo);
                    break;
                case 3:
                    JSONObject footer = JSONObject.parseObject(setting.getValue());
                    SettingVo.FooterVo footerVo = footer.toJavaObject(SettingVo.FooterVo.class);
                    footerVo.setOrderNum(setting.getOrderNum());
                    footerList.add(footerVo);
                    break;
            }
        });
        SettingVo settingVo = JSONObject.parseObject(JSONObject.toJSONString(infoMap), SettingVo.class);
        settingVo.setBadge(badgeList);
        settingVo.setFooter(footerList);

        return settingVo;
    }

    @Log
    @Override
    public boolean updateBaseSetting(SettingDto settingDto) {
        // update all properties
        ObjectMapper objectMapper = new ObjectMapper();
        Map<String, Object> map = objectMapper.convertValue(settingDto, Map.class);

        map.forEach((key, value) -> {
            if (value == null) {
                return;
            }
            switch (key) {
                case "badge":
                    List<SettingDto.BadgeDto> badgeList = JSONObject.parseObject(
                            JSONObject.toJSONString(value), new TypeReference<List<SettingDto.BadgeDto>>() {
                            });
                    badgeList.forEach(badge -> {
                        Setting setting = getOne(
                                new LambdaQueryWrapper<Setting>().eq(Setting::getName, "badge_" + badge.getIndex()));
                        if (setting == null) {
                            setting = new Setting();
                            setting.setName("badge_" + badge.getIndex());
                            setting.setValue(JSONObject.toJSONString(badge));
                            setting.setType(2);
                            setting.setStatus(SystemConstants.STATUS_NORMAL);
                            setting.setOrderNum(badge.getOrderNum());
                            save(setting);
                        } else {
                            setting.setValue(JSONObject.toJSONString(badge));
                            updateById(setting);
                        }
                    });
                    break;
                case "footer":
                    List<SettingDto.FooterDto> footerList = JSONObject.parseObject(
                            JSONObject.toJSONString(value), new TypeReference<List<SettingDto.FooterDto>>() {
                            });
                    footerList.forEach(footer -> {
                        Setting setting = getOne(
                                new LambdaQueryWrapper<Setting>().eq(Setting::getName, "footer_" + footer.getIndex()));
                        if (setting == null) {
                            setting = new Setting();
                            setting.setName("footer_" + footer.getIndex());
                            setting.setValue(JSONObject.toJSONString(footer));
                            setting.setType(3);
                            setting.setStatus(SystemConstants.STATUS_NORMAL);
                            setting.setOrderNum(footer.getOrderNum());
                            save(setting);
                        } else {
                            setting.setValue(JSONObject.toJSONString(footer));
                            updateById(setting);
                        }
                    });
                    break;

                default:
                    LambdaQueryWrapper<Setting> wrapper = new LambdaQueryWrapper<>();
                    wrapper.eq(Setting::getName, key);
                    Setting setting = getOne(wrapper);
                    if (setting == null) {
                        setting = new Setting();
                        setting.setName(key);
                        setting.setValue(value.toString());
                        setting.setType(1);
                        setting.setStatus(SystemConstants.STATUS_NORMAL);
                        setting.setOrderNum(0);
                        save(setting);
                    } else {
                        setting.setValue(value.toString());
                        updateById(setting);
                    }
                    break;
            }
        });
        return true;
    }
}
