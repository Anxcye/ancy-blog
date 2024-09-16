package com.anxcye.utils;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.util.Date;
import java.util.UUID;

/**
 * JWT工具类
 */
public class JwtUtil {

    // 有效期为
    public static final Long JWT_TTL = 24 * 60 * 60 * 1000L;// 60 * 60 *1000 一个小时
    // 设置秘钥明文
    public static final String JWT_KEY = "anxJWT signature does not match locally computed signature. JWT validity cannot be asserted and should not be trusted.cye";

    public static String getUUID() {
        return UUID.randomUUID().toString().replaceAll("-", "");
    }

    /**
     * 生成jtw
     *
     * @param subject token中要存放的数据（json格式）
     */
    public static String createJWT(String subject) {
        return getJwtBuilder(subject, null, getUUID());
    }

    /**
     * 生成jtw
     *
     * @param subject   token中要存放的数据（json格式）
     * @param ttlMillis token超时时间
     */
    public static String createJWT(String subject, Long ttlMillis) {
        return getJwtBuilder(subject, ttlMillis, getUUID());// 设置过期时间
    }

    private static String getJwtBuilder(String subject, Long ttlMillis, String uuid) {
        long nowMillis = System.currentTimeMillis();
        Date now = new Date(nowMillis);
        if (ttlMillis == null) {
            ttlMillis = JwtUtil.JWT_TTL;
        }
        long expMillis = nowMillis + ttlMillis;
        Date expDate = new Date(expMillis);
        return Jwts.builder()
                .subject(subject)
                // .issuer("ancy")
                .issuedAt(now)
                .signWith(generalKey())
                .expiration(expDate)
                .compact();
    }

    /**
     * 创建token
     */
    public static String createJWT(String id, String subject, Long ttlMillis) {
        return getJwtBuilder(subject, ttlMillis, id);// 设置过期时间
    }

    /**
     * 生成加密后的秘钥 secretKey
     */

    public static SecretKey generalKey() {
        byte[] encodedKey = JWT_KEY.getBytes();
        return new SecretKeySpec(encodedKey, 0, encodedKey.length, "HmacSHA256");
    }

    /**
     * 解析
     */
    public static Claims parseJWT(String jwt) throws Exception {
        return Jwts.parser()
                .verifyWith(generalKey())
                .build()
                .parseSignedClaims(jwt)
                .getPayload();
    }
}