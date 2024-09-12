package com.anxcye.utils;

import org.apache.poi.ss.formula.functions.T;
import org.springframework.beans.BeanUtils;

import java.util.List;
import java.util.stream.Collectors;

public class BeanCopyUtils {

    private BeanCopyUtils() {
    }

    public static <T> T copyBean(Object source, Class<T> target) {
        if (source == null) {
            return null;
        }
        try {
            T t = target.newInstance();
            BeanUtils.copyProperties(source, t);
            return t;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public static <T> List<T> copyList(List<?> source, Class<T> target) {
        if (source == null || source.isEmpty()) {
            return List.of();
        }
        try {
            return source.stream()
                    .map(item -> copyBean(item, target))
                    .collect(Collectors.toList());
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

}
